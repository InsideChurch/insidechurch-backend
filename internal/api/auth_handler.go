package api

import (
	"encoding/json"
	"net/http"

	"insidechurch.com/backend/internal/models"
	"insidechurch.com/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenString, isGlobalSuperAdmin, userEmail, userName, userRole, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := models.LoginResponse{
		Token:              tokenString,
		IsGlobalSuperAdmin: isGlobalSuperAdmin,
		UserEmail:          userEmail,
		UserName:           userName,
		UserRole:           userRole,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) CreateTenantSuperAdmin(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTenantSuperAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.CreateTenantSuperAdmin(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}