package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"insidechurch.com/backend/internal/models"     // Adjust module path
	"insidechurch.com/backend/internal/repository" // Adjust module path
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
