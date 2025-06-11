package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	userContextKey contextKey = "user"
)

type AuthClaims struct {
	UserID             uuid.UUID  `json:"user_id"`
	Email              string     `json:"email"`
	IsGlobalSuperAdmin bool       `json:"is_global_super_admin"`
	TenantID           *uuid.UUID `json:"tenant_id,omitempty"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("JWT_SECRET environment variable not set. Cannot authenticate requests.")
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := &AuthClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			log.Printf("Token parsing error: %v", err)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (*AuthClaims, error) {
	claims, ok := ctx.Value(userContextKey).(*AuthClaims)
	if !ok {
		return nil, errors.New("user claims not found in context")
	}
	return claims, nil
}

func GlobalAdminRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := GetUserFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized: No user info", http.StatusUnauthorized)
			return
		}
		if !claims.IsGlobalSuperAdmin {
			http.Error(w, "Forbidden: Global Super Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
