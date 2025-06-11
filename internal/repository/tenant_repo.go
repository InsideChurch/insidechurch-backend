package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"insidechurch.com/backend/internal/models"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) CreateTenant(tenant *models.Tenant) error {
	tenant.ID = uuid.New()
	query := `INSERT INTO tenants (id, name, type, parent_id) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, tenant.ID, tenant.Name, tenant.Type, tenant.ParentID)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	return nil
}
