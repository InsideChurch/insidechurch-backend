package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"insidechurch.com/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, email, password_hash, name, tenant_id, is_global_super_admin, created_at, updated_at FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.TenantID,
		&user.IsGlobalSuperAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	user.ID = uuid.New()
	query := `INSERT INTO users (id, email, password_hash, name, tenant_id, is_global_super_admin)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash, user.Name, user.TenantID, user.IsGlobalSuperAdmin)
	return err
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	var tenantID sql.NullString
	var isGlobalSuperAdmin sql.NullBool

	query := "SELECT id, email, password_hash, name, role, tenant_id, is_global_super_admin, created_at, updated_at FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).
		Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&user.Role,
			&tenantID,
			&isGlobalSuperAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if tenantID.Valid {
		parsedUUID, err := uuid.Parse(tenantID.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tenant_id UUID from DB: %w", err)
		}
		user.TenantID = &parsedUUID
	} else {
		user.TenantID = nil
	}

	if isGlobalSuperAdmin.Valid {
		user.IsGlobalSuperAdmin = isGlobalSuperAdmin.Bool
	} else {
		user.IsGlobalSuperAdmin = false
	}

	return &user, nil
}

func (r *UserRepository) CreateUserWithTenantAndRole(user *models.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `INSERT INTO users (id, email, password_hash, name, role, tenant_id, is_global_super_admin, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Role,
		user.TenantID,
		user.IsGlobalSuperAdmin,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user with tenant and role: %w", err)
	}
	return nil
}
