package repository

import (
	"database/sql"

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
