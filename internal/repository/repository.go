package repository

import (
	"github.com/ZhangPengPaul/LightOTA/internal/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(
		&model.Tenant{},
		&model.Product{},
		&model.Device{},
		&model.Firmware{},
		&model.UpgradeTask{},
		&model.DeviceUpgradeRecord{},
	)
}

type ITenantRepository interface {
	FindByAPIKey(apiKey string) (*model.Tenant, error)
	FindByID(id string) (*model.Tenant, error)
	List(limit, offset int) ([]model.Tenant, int64, error)
	Create(tenant *model.Tenant) error
	Update(tenant *model.Tenant) error
}

type IProductRepository interface {
	FindByID(tenantID, id string) (*model.Product, error)
	ListByTenant(tenantID string, limit, offset int) ([]model.Product, int64, error)
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(product *model.Product) error
}

type IFirmwareRepository interface {
	FindByID(tenantID, id string) (*model.Firmware, error)
	ListByProduct(tenantID, productID string, limit, offset int) ([]model.Firmware, int64, error)
	Create(firmware *model.Firmware) error
	Delete(firmware *model.Firmware) error
	GetLatestActive(tenantID, productID string) (*model.Firmware, error)
}

type IDeviceRepository interface {
	FindByID(id string) (*model.Device, error)
	Create(device *model.Device) error
	UpdateCurrentVersion(id string, version string) error
}

type IUpgradeTaskRepository interface {
	FindByID(tenantID, id string) (*model.UpgradeTask, error)
	ListByTenant(tenantID string, limit, offset int) ([]model.UpgradeTask, int64, error)
	ListByProduct(tenantID, productID string, limit, offset int) ([]model.UpgradeTask, int64, error)
	Create(task *model.UpgradeTask) error
	Update(task *model.UpgradeTask) error
	CountByStatus(taskID string) (map[model.DeviceUpgradeStatus]int, error)
}

type IDeviceUpgradeRecordRepository interface {
	Create(record *model.DeviceUpgradeRecord) error
	UpdateStatus(id string, status model.DeviceUpgradeStatus) error
	UpdateStatusWithError(id string, status model.DeviceUpgradeStatus, errorMsg string) error
	FindByTaskAndDevice(taskID, deviceID string) (*model.DeviceUpgradeRecord, error)
	ListPending(taskID string, limit int) ([]model.DeviceUpgradeRecord, error)
}

func (r *Repository) FindByAPIKey(apiKey string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.Where("api_key = ?", apiKey).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) FindByID(id string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.Where("id = ?", id).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) List(limit, offset int) ([]model.Tenant, int64, error) {
	var tenants []model.Tenant
	var total int64
	err := r.db.Model(&model.Tenant{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Offset(offset).Limit(limit).Find(&tenants).Error
	return tenants, total, err
}

func (r *Repository) Create(tenant *model.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *Repository) Update(tenant *model.Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *Repository) FindProductByID(tenantID, id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) ListProductsByTenant(tenantID string, limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64
	err := r.db.Model(&model.Product{}).Where("tenant_id = ?", tenantID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Where("tenant_id = ?", tenantID).Offset(offset).Limit(limit).Find(&products).Error
	return products, total, err
}

func (r *Repository) CreateProduct(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *Repository) UpdateProduct(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *Repository) DeleteProduct(product *model.Product) error {
	return r.db.Delete(product).Error
}

func (r *Repository) FindFirmwareByID(tenantID, id string) (*model.Firmware, error) {
	var firmware model.Firmware
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&firmware).Error
	if err != nil {
		return nil, err
	}
	return &firmware, nil
}

func (r *Repository) ListFirmwaresByProduct(tenantID, productID string, limit, offset int) ([]model.Firmware, int64, error) {
	var firmwares []model.Firmware
	var total int64
	err := r.db.Model(&model.Firmware{}).Where("tenant_id = ? AND product_id = ?", tenantID, productID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Where("tenant_id = ? AND product_id = ?", tenantID, productID).Order("version_code DESC").Offset(offset).Limit(limit).Find(&firmwares).Error
	return firmwares, total, err
}

func (r *Repository) CreateFirmware(firmware *model.Firmware) error {
	return r.db.Create(firmware).Error
}

func (r *Repository) DeleteFirmware(firmware *model.Firmware) error {
	return r.db.Delete(firmware).Error
}

func (r *Repository) Delete(firmware *model.Firmware) error {
	return r.DeleteFirmware(firmware)
}

func (r *Repository) GetLatestActive(tenantID, productID string) (*model.Firmware, error) {
	var firmware model.Firmware
	err := r.db.Where("tenant_id = ? AND product_id = ? AND is_active = ?", tenantID, productID, true).Order("version_code DESC").First(&firmware).Error
	if err != nil {
		return nil, err
	}
	return &firmware, nil
}

func (r *Repository) GetLatestActiveFirmware(tenantID, productID string) (*model.Firmware, error) {
	return r.GetLatestActive(tenantID, productID)
}

func (r *Repository) FindDeviceByID(id string) (*model.Device, error) {
	var device model.Device
	err := r.db.Where("id = ?", id).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *Repository) CreateDevice(device *model.Device) error {
	return r.db.Create(device).Error
}

func (r *Repository) FindDeviceByExternalID(tenantID, externalID string) (*model.Device, error) {
	var device model.Device
	err := r.db.Where("tenant_id = ? AND external_device_id = ?", tenantID, externalID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *Repository) UpdateDeviceCurrentVersion(id string, version string) error {
	return r.db.Model(&model.Device{}).Where("id = ?", id).Update("current_version", version).Error
}

func (r *Repository) FindUpgradeTaskByID(tenantID, id string) (*model.UpgradeTask, error) {
	var task model.UpgradeTask
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) ListUpgradeTasksByTenant(tenantID string, limit, offset int) ([]model.UpgradeTask, int64, error) {
	var tasks []model.UpgradeTask
	var total int64
	err := r.db.Model(&model.UpgradeTask{}).Where("tenant_id = ?", tenantID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Where("tenant_id = ?", tenantID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&tasks).Error
	return tasks, total, err
}

func (r *Repository) ListUpgradeTasksByProduct(tenantID, productID string, limit, offset int) ([]model.UpgradeTask, int64, error) {
	var tasks []model.UpgradeTask
	var total int64
	err := r.db.Model(&model.UpgradeTask{}).Where("tenant_id = ? AND product_id = ?", tenantID, productID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Where("tenant_id = ? AND product_id = ?", tenantID, productID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&tasks).Error
	return tasks, total, err
}

func (r *Repository) CreateUpgradeTask(task *model.UpgradeTask) error {
	return r.db.Create(task).Error
}

func (r *Repository) UpdateUpgradeTask(task *model.UpgradeTask) error {
	return r.db.Save(task).Error
}

func (r *Repository) CountRecordsByStatus(taskID string) (map[model.DeviceUpgradeStatus]int, error) {
	result := make(map[model.DeviceUpgradeStatus]int)
	rows, err := r.db.Table("device_upgrade_records").Where("task_id = ?", taskID).Select("status, count(*) as count").Group("status").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var status model.DeviceUpgradeStatus
		var count int
		rows.Scan(&status, &count)
		result[status] = count
	}
	return result, err
}

func (r *Repository) CreateDeviceUpgradeRecord(record *model.DeviceUpgradeRecord) error {
	return r.db.Create(record).Error
}

func (r *Repository) UpdateDeviceUpgradeRecordStatus(id string, status model.DeviceUpgradeStatus) error {
	return r.db.Model(&model.DeviceUpgradeRecord{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Repository) UpdateDeviceUpgradeRecordStatusWithError(id string, status model.DeviceUpgradeStatus, errorMsg string) error {
	return r.db.Model(&model.DeviceUpgradeRecord{}).Where("id = ?", id).Updates(map[string]interface{}{"status": status, "error_message": errorMsg}).Error
}

func (r *Repository) FindDeviceUpgradeRecordByTaskAndDevice(taskID, deviceID string) (*model.DeviceUpgradeRecord, error) {
	var record model.DeviceUpgradeRecord
	err := r.db.Where("task_id = ? AND device_id = ?", taskID, deviceID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *Repository) ListPendingDeviceUpgradeRecords(taskID string, limit int) ([]model.DeviceUpgradeRecord, error) {
	var records []model.DeviceUpgradeRecord
	err := r.db.Where("task_id = ? AND status = ?", taskID, model.DeviceStatusPending).Limit(limit).Find(&records).Error
	return records, err
}
