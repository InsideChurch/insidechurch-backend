package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID  `json:"id"`
	Email              string     `json:"email"`
	PasswordHash       string     `json:"-"`
	Name               string     `json:"name"`
	Role               string     `json:"role"`
	TenantID           *uuid.UUID `json:"tenant_id,omitempty"`
	IsGlobalSuperAdmin bool       `json:"is_global_super_admin"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token          string `json:"token"`
	IsGlobalSuperAdmin bool `json:"is_global_super_admin"`
	UserEmail      string `json:"user_email"` 
	UserName       string `json:"user_name"`  
	UserRole       string `json:"user_role"`
}

type CreateTenantSuperAdminRequest struct {
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	TenantID uuid.UUID `json:"tenant_id"`
}