package service

import (
	"github.com/ZhangPengPaul/LightOTA/internal/model"
	"github.com/ZhangPengPaul/LightOTA/internal/repository"
	"github.com/google/uuid"
)

type TenantService struct {
	repo *repository.Repository
}

func NewTenantService(repo *repository.Repository) *TenantService {
	return &TenantService{repo: repo}
}

type CreateTenantRequest struct {
	Name              string `json:"name" binding:"required"`
	ExternalDeviceAPIURL string `json:"external_device_api_url"`
}

func (s *TenantService) Create(req *CreateTenantRequest) (*model.Tenant, error) {
	tenant := &model.Tenant{
		ID:                uuid.New().String(),
		Name:              req.Name,
		APIKey:            uuid.New().String(),
		ExternalDeviceAPIURL: req.ExternalDeviceAPIURL,
	}
	err := s.repo.Create(tenant)
	return tenant, err
}

func (s *TenantService) List(limit, offset int) ([]model.Tenant, int64, error) {
	return s.repo.List(limit, offset)
}

func (s *TenantService) GetByID(id string) (*model.Tenant, error) {
	return s.repo.FindByID(id)
}

type UpdateTenantRequest struct {
	Name              *string `json:"name"`
	ExternalDeviceAPIURL *string `json:"external_device_api_url"`
}

func (s *TenantService) Update(id string, req *UpdateTenantRequest) (*model.Tenant, error) {
	tenant, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.ExternalDeviceAPIURL != nil {
		tenant.ExternalDeviceAPIURL = *req.ExternalDeviceAPIURL
	}
	err = s.repo.Update(tenant)
	return tenant, err
}
