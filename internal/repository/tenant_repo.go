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

func (r *TenantRepository) GetAllTenants() ([]models.Tenant, error) {
	query := "SELECT id, name, type, parent_id, created_at, updated_at FROM tenants ORDER BY created_at DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tenants: %w", err)
	}
	defer rows.Close()

	var tenants []models.Tenant
	for rows.Next() {
		var tenant models.Tenant

		err := rows.Scan(&tenant.ID, &tenant.Name, &tenant.Type, &tenant.ParentID, &tenant.CreatedAt, &tenant.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tenant row: %w", err)
		}
		tenants = append(tenants, tenant)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return tenants, nil
}