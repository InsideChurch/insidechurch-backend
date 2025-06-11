package service

import (
	"errors"
	"fmt"

	"insidechurch.com/backend/internal/models"
	"insidechurch.com/backend/internal/repository"
)

type TenantService struct {
	tenantRepo *repository.TenantRepository
}

func NewTenantService(tenantRepo *repository.TenantRepository) *TenantService {
	return &TenantService{tenantRepo: tenantRepo}
}

func (s *TenantService) CreateTenant(req *models.CreateTenantRequest) (*models.Tenant, error) {
	if req.Name == "" || req.Type == "" {
		return nil, errors.New("tenant name and type are required")
	}

	tenant := &models.Tenant{
		Name:     req.Name,
		Type:     req.Type,
		ParentID: req.ParentID,
	}

	err := s.tenantRepo.CreateTenant(tenant)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create tenant: %w", err)
	}
	return tenant, nil
}

func (s *TenantService) GetAllTenants() ([]models.TenantResponse, error) {
	tenants, err := s.tenantRepo.GetAllTenants()
	if err != nil {
		return nil, fmt.Errorf("service: failed to get all tenants: %w", err)
	}
	return tenants, nil
}