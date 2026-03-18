package service

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ZhangPengPaul/LightOTA/internal/config"
	"github.com/ZhangPengPaul/LightOTA/internal/httpsse"
	"github.com/ZhangPengPaul/LightOTA/internal/model"
	"github.com/ZhangPengPaul/LightOTA/internal/mqtt"
	"github.com/ZhangPengPaul/LightOTA/internal/repository"
	"github.com/google/uuid"
)

type UpgradeService struct {
	repo         *repository.Repository
	mqttClient   *mqtt.Client
	httpSSE      *httpsse.Manager
	cfg          *config.Config
}

func NewUpgradeService(
	taskRepo *repository.Repository,
	recordRepo *repository.Repository,
	deviceRepo *repository.Repository,
	firmwareRepo *repository.Repository,
	mqttClient *mqtt.Client,
	cfg *config.Config,
) *UpgradeService {
	return &UpgradeService{
		repo:         taskRepo,
		mqttClient:   mqttClient,
		httpSSE:      httpsse.NewManager(),
		cfg:          cfg,
	}
}

type CreateUpgradeTaskRequest struct {
	ProductID       string             `json:"product_id" binding:"required"`
	FirmwareID       string             `json:"firmware_id" binding:"required"`
	TaskName         string             `json:"task_name" binding:"required"`
	UpgradeType      model.UpgradeType  `json:"upgrade_type" binding:"required"`
	GrayPercent      int                `json:"gray_percent"`
	TargetDeviceIDs  []string           `json:"target_device_ids"`
	PushRate         int                `json:"push_rate"`
}

func (s *UpgradeService) CreateTask(tenantID string, req *CreateUpgradeTaskRequest, tenant *model.Tenant) (*model.UpgradeTask, error) {
	task := &model.UpgradeTask{
		ID:             uuid.New().String(),
		TenantID:       tenantID,
		ProductID:      req.ProductID,
		FirmwareID:     req.FirmwareID,
		TaskName:       req.TaskName,
		UpgradeType:    req.UpgradeType,
		GrayPercent:    req.GrayPercent,
		PushRate:       req.PushRate,
		Status:         model.TaskStatusCreated,
		CreatedBy:      "",
	}

	if req.PushRate <= 0 {
		task.PushRate = 10
	}

	var devices []model.ThirdPartyDevice

	switch req.UpgradeType {
	case model.UpgradeTypeSpecified:
		for _, deviceID := range req.TargetDeviceIDs {
			devices = append(devices, model.ThirdPartyDevice{
				DeviceID: deviceID,
			})
		}
	case model.UpgradeTypeAll:
		if tenant.ExternalDeviceAPIURL != "" {
			queryReq := &model.ThirdPartyDeviceQueryRequest{
				ProductID:       req.ProductID,
				Percent:         100,
				ExcludeVersions: []string{},
				OnlyOnline:      true,
			}
			client := NewThirdPartyClient(tenant.ExternalDeviceAPIURL, tenant.APIKey)
			var err error
			devices, _, err = client.QueryDevices(queryReq)
			if err != nil {
				return nil, err
			}
		}
	case model.UpgradeTypeGray:
		if tenant.ExternalDeviceAPIURL != "" {
			queryReq := &model.ThirdPartyDeviceQueryRequest{
				ProductID:       req.ProductID,
				Percent:         req.GrayPercent,
				ExcludeVersions: []string{},
				OnlyOnline:      true,
			}
			client := NewThirdPartyClient(tenant.ExternalDeviceAPIURL, tenant.APIKey)
			var err error
			devices, _, err = client.QueryDevices(queryReq)
			if err != nil {
				return nil, err
			}
		}
	}

	task.TargetDevicesCount = len(devices)

	err := s.repo.CreateUpgradeTask(task)
	if err != nil {
		return nil, err
	}

	for _, dev := range devices {
		device, err := s.repo.FindDeviceByID(dev.DeviceID)
		if err != nil {
			newDevice := &model.Device{
				ID:               uuid.New().String(),
				TenantID:         tenantID,
				ProductID:        req.ProductID,
				ExternalDeviceID: dev.DeviceID,
				CurrentVersion:   dev.CurrentVersion,
			}
			err = s.repo.CreateDevice(newDevice)
			if err != nil {
				continue
			}
			device = newDevice
		}

		record := &model.DeviceUpgradeRecord{
			ID:         uuid.New().String(),
			TaskID:     task.ID,
			DeviceID:   device.ID,
			OldVersion: device.CurrentVersion,
			NewVersion: "",
			Status:     model.DeviceStatusPending,
		}
		s.repo.CreateDeviceUpgradeRecord(record)
	}

	return task, nil
}

func (s *UpgradeService) ListTasks(tenantID string, productID *string, limit, offset int) ([]model.UpgradeTask, int64, error) {
	if productID != nil && *productID != "" {
		return s.repo.ListUpgradeTasksByProduct(tenantID, *productID, limit, offset)
	}
	return s.repo.ListUpgradeTasksByTenant(tenantID, limit, offset)
}

func (s *UpgradeService) GetTask(tenantID, id string) (*model.UpgradeTask, map[model.DeviceUpgradeStatus]int, error) {
	task, err := s.repo.FindUpgradeTaskByID(tenantID, id)
	if err != nil {
		return nil, nil, err
	}

	counts, err := s.repo.CountRecordsByStatus(task.ID)
	return task, counts, err
}

func (s *UpgradeService) UpdateDeviceRecordStatus(recordID string, status model.DeviceUpgradeStatus, errorMsg string) error {
	if errorMsg != "" {
		return s.repo.UpdateDeviceUpgradeRecordStatusWithError(recordID, status, errorMsg)
	}
	return s.repo.UpdateDeviceUpgradeRecordStatus(recordID, status)
}

func (s *UpgradeService) FindPendingRecords(taskID string, limit int) ([]model.DeviceUpgradeRecord, error) {
	return s.repo.ListPendingDeviceUpgradeRecords(taskID, limit)
}

func (s *UpgradeService) PushNotification(record *model.DeviceUpgradeRecord, firmware *model.Firmware, device *model.Device) error {
	if s.mqttClient != nil && s.cfg.MQTT.Enabled {
		err := s.mqttClient.PublishUpgradeNotification(device.ExternalDeviceID, firmware)
		if err != nil {
			return err
		}
	}

	if s.httpSSE != nil {
		s.httpSSE.NotifyUpgrade(device.ExternalDeviceID, firmware)
	}

	return s.repo.UpdateDeviceUpgradeRecordStatus(record.ID, model.DeviceStatusNotified)
}

func (s *UpgradeService) GetStats(taskID string) (map[model.DeviceUpgradeStatus]int, error) {
	return s.repo.CountRecordsByStatus(taskID)
}

func (s *UpgradeService) StartTask(task *model.UpgradeTask) error {
	now := time.Now()
	task.Status = model.TaskStatusRunning
	task.StartedAt = &now
	return s.repo.UpdateUpgradeTask(task)
}

func (s *UpgradeService) CompleteTask(task *model.UpgradeTask) error {
	now := time.Now()
	task.Status = model.TaskStatusCompleted
	task.CompletedAt = &now
	return s.repo.UpdateUpgradeTask(task)
}

func (s *UpgradeService) SelectRandomDevices(devices []model.ThirdPartyDevice, percent int) []model.ThirdPartyDevice {
	if percent >= 100 {
		return devices
	}
	targetCount := len(devices) * percent / 100
	if targetCount <= 0 {
		targetCount = 1
	}

	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(devices))
	result := make([]model.ThirdPartyDevice, 0, targetCount)
	for i := 0; i < targetCount; i++ {
		result = append(result, devices[perm[i]])
	}
	return result
}

func (s *UpgradeService) GetFirmwareByID(tenantID, id string) (*model.Firmware, error) {
	return s.repo.FindFirmwareByID(tenantID, id)
}

func (s *UpgradeService) GetDeviceByID(id string) (*model.Device, error) {
	return s.repo.FindDeviceByID(id)
}



func (s *UpgradeService) GetLatestActiveFirmware(tenantID, productID string) (*model.Firmware, error) {
	return s.repo.GetLatestActive(tenantID, productID)
}

func (s *UpgradeService) GetFirmwareFile(tenantID, id string) (*os.File, int64, error) {
	firmware, err := s.repo.FindFirmwareByID(tenantID, id)
	if err != nil {
		return nil, 0, err
	}
	file, err := os.Open(firmware.FilePath)
	if err != nil {
		return nil, 0, err
	}
	return file, firmware.FileSize, nil
}

func (s *UpgradeService) FindPendingByTaskAndDevice(taskID, deviceID string) (*model.DeviceUpgradeRecord, error) {
	return s.repo.FindDeviceUpgradeRecordByTaskAndDevice(taskID, deviceID)
}

func (s *UpgradeService) HandleSSE(w http.ResponseWriter, r *http.Request, deviceID string) {
	s.httpSSE.HandleSSE(w, r, deviceID)
}
