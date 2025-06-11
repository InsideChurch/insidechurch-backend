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

func (r *TenantRepository) GetAllTenants() ([]models.TenantResponse, error) {
	query := `
	    SELECT
	        t.id, t.name, t.type, t.parent_id, t.created_at, t.updated_at,
	        p.name AS parent_name -- Select parent's name with an alias
	    FROM
	        tenants t
	    LEFT JOIN
	        tenants p ON t.parent_id = p.id -- LEFT JOIN to get parent details
	    ORDER BY
	        t.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tenants: %w", err)
	}
	defer rows.Close()

	var tenants []models.TenantResponse
	for rows.Next() {
		var tenant models.TenantResponse
		var parentID sql.NullString 
		var parentName sql.NullString

		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Type,
			&parentID,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
			&parentName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tenant row: %w", err)
		}

		if parentID.Valid {
			parsedUUID, err := uuid.Parse(parentID.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse parent_id UUID: %w", err)
			}
			tenant.ParentID = &parsedUUID
		} else {
			tenant.ParentID = nil
		}

		if parentName.Valid {
			tenant.ParentName = &parentName.String
		} else {
			tenant.ParentName = nil
		}

		tenants = append(tenants, tenant)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return tenants, nil
}