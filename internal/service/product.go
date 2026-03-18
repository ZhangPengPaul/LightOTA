package service

import (
	"github.com/ZhangPengPaul/LightOTA/internal/model"
	"github.com/ZhangPengPaul/LightOTA/internal/repository"
	"github.com/google/uuid"
)

type ProductService struct {
	repo *repository.Repository
}

func NewProductService(repo *repository.Repository) *ProductService {
	return &ProductService{repo: repo}
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *ProductService) Create(tenantID string, req *CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
	}
	err := s.repo.CreateProduct(product)
	return product, err
}

func (s *ProductService) List(tenantID string, limit, offset int) ([]model.Product, int64, error) {
	return s.repo.ListProductsByTenant(tenantID, limit, offset)
}

func (s *ProductService) GetByID(tenantID, id string) (*model.Product, error) {
	return s.repo.FindProductByID(tenantID, id)
}

type UpdateProductRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (s *ProductService) Update(tenantID, id string, req *UpdateProductRequest) (*model.Product, error) {
	product, err := s.repo.FindProductByID(tenantID, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	err = s.repo.UpdateProduct(product)
	return product, err
}

func (s *ProductService) Delete(tenantID, id string) error {
	product, err := s.repo.FindProductByID(tenantID, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteProduct(product)
}
