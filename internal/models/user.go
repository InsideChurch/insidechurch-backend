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
	Token string `json:"token"`
	User  User   `json:"user"`
}
