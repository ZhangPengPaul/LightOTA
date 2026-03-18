package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ZhangPengPaul/LightOTA/internal/model"
	"github.com/ZhangPengPaul/LightOTA/internal/repository"
	"github.com/google/uuid"
)

type FirmwareService struct {
	repo         *repository.Repository
	storagePath  string
}

func NewFirmwareService(repo *repository.Repository, storagePath string) *FirmwareService {
	os.MkdirAll(storagePath, 0755)
	return &FirmwareService{repo: repo, storagePath: storagePath}
}

type CreateFirmwareRequest struct {
	ProductID     string `form:"product_id" json:"product_id" binding:"required"`
	Version       string `form:"version" json:"version" binding:"required"`
	VersionCode   int    `form:"versionCode" json:"version_code" binding:"required"`
	Changelog     string `form:"changelog" json:"changelog"`
	ReleaseNotes  string `form:"releaseNotes" json:"release_notes"`
}

func (s *FirmwareService) Create(tenantID string, req *CreateFirmwareRequest, fileContent io.Reader, filename string) (*model.Firmware, error) {
	_, err := s.repo.FindProductByID(tenantID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	firmwareID := uuid.New().String()
	ext := filepath.Ext(filename)
	savePath := filepath.Join(s.storagePath, fmt.Sprintf("%s%s", firmwareID, ext))

	out, err := os.Create(savePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	hash := md5.New()
	multiWriter := io.MultiWriter(out, hash)
	fileSize, err := io.Copy(multiWriter, fileContent)
	if err != nil {
		return nil, err
	}

	md5sum := hex.EncodeToString(hash.Sum(nil))

	firmware := &model.Firmware{
		ID:          firmwareID,
		TenantID:    tenantID,
		ProductID:   req.ProductID,
		Version:     req.Version,
		VersionCode: req.VersionCode,
		Changelog:   req.Changelog,
		FilePath:    savePath,
		FileSize:    fileSize,
		MD5:         md5sum,
		ReleaseNotes: req.ReleaseNotes,
		IsActive:    true,
	}

	err = s.repo.CreateFirmware(firmware)
	return firmware, err
}

func (s *FirmwareService) List(tenantID, productID string, limit, offset int) ([]model.Firmware, int64, error) {
	return s.repo.ListFirmwaresByProduct(tenantID, productID, limit, offset)
}

func (s *FirmwareService) GetByID(tenantID, id string) (*model.Firmware, error) {
	return s.repo.FindFirmwareByID(tenantID, id)
}

func (s *FirmwareService) Delete(tenantID, id string) error {
	firmware, err := s.repo.FindFirmwareByID(tenantID, id)
	if err != nil {
		return err
	}
	if err := os.Remove(firmware.FilePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return s.repo.DeleteFirmware(firmware)
}

func (s *FirmwareService) GetFile(tenantID, id string) (*os.File, int64, error) {
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

func (s *FirmwareService) GetLatestActive(tenantID, productID string) (*model.Firmware, error) {
	return s.repo.GetLatestActiveFirmware(tenantID, productID)
}
