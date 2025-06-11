package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"insidechurch.com/backend/internal/models"     
	"insidechurch.com/backend/internal/repository" 
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) AuthenticateUser(email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id":               user.ID,
		"email":                 user.Email,
		"is_global_super_admin": user.IsGlobalSuperAdmin,
		"tenant_id":             user.TenantID,
		"exp":                   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) Login(email, password string) (string, bool, string, string, string, error) { 
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return "", false, "", "", "", fmt.Errorf("service: failed to find user for login: %w", err)
	}
	if user == nil {
		return "", false, "", "", "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", false, "", "", "", errors.New("invalid credentials")
	}

	isGlobalSuperAdmin := user.IsGlobalSuperAdmin 

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"name":     user.Name,
		"role":     user.Role,
		"tenant_id": user.TenantID,
		"is_global_super_admin": isGlobalSuperAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", false, "", "", "", fmt.Errorf("service: failed to sign token: %w", err)
	}

	return tokenString, isGlobalSuperAdmin, user.Email, user.Name, user.Role, nil 
}

func (s *AuthService) CreateTenantSuperAdmin(req *models.CreateTenantSuperAdminRequest) (*models.User, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" || req.TenantID == uuid.Nil { // <-- uuid.Nil is now recognized
		return nil, errors.New("email, password, name, and tenant ID are required")
	}

	existingUser, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("service: failed to check for existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("service: failed to hash password: %w", err)
	}

	user := &models.User{
		Email:    req.Email,
		PasswordHash: string(hashedPassword), 
		Name:     req.Name,
		Role:     "tenant_super_admin",
		TenantID: &req.TenantID,
        IsGlobalSuperAdmin: false,
	}

	if err := s.userRepo.CreateUserWithTenantAndRole(user); err != nil {
		return nil, fmt.Errorf("service: failed to create tenant super admin in repository: %w", err)
	}

	return user, nil
}