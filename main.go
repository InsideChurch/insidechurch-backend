package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"insidechurch.com/backend/internal/api"
	"insidechurch.com/backend/internal/repository"
	"insidechurch.com/backend/internal/service"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set. Please set it in your .env file or environment.")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL database!")

	schemaSQL := `
	    CREATE TABLE IF NOT EXISTS tenants (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(255) NOT NULL,
            type VARCHAR(50) NOT NULL,
            parent_id UUID REFERENCES tenants(id) NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            name VARCHAR(255) NOT NULL,
            role VARCHAR(50) NOT NULL DEFAULT 'tenant_admin',
            tenant_id UUID REFERENCES tenants(id) NULL,
            is_global_super_admin BOOLEAN DEFAULT FALSE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE IF NOT EXISTS roles (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(50) UNIQUE NOT NULL
        );
        CREATE TABLE IF NOT EXISTS user_roles (
            user_id UUID REFERENCES users(id),
            role_id UUID REFERENCES roles(id),
            tenant_id UUID REFERENCES tenants(id),
            PRIMARY KEY (user_id, role_id, tenant_id)
        );
        CREATE TABLE IF NOT EXISTS members (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            tenant_id UUID NOT NULL REFERENCES tenants(id),
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE NULL,
            phone_number VARCHAR(50) UNIQUE NULL,
            birthday DATE NOT NULL,
            address TEXT NULL,
            membership_status VARCHAR(50) NOT NULL DEFAULT 'Active',
            marital_status VARCHAR(50) NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
        INSERT INTO roles (name) VALUES ('tenant_super_admin'), ('tenant_admin'), ('leadership');
	`
	_, err = db.Exec(schemaSQL)
	if err != nil {
		log.Printf("Warning: Error creating database schema: %v (This might be okay if tables/roles already exist)", err)
	} else {
		fmt.Println("Database schema initialized successfully.")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to InsideChurch Backend MVP!")
}

func main() {
	initDB()

	r := mux.NewRouter()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set. Please set it in your .env file or environment.")
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := api.NewAuthHandler(authService)

	tenantRepo := repository.NewTenantRepository(db)
	tenantService := service.NewTenantService(tenantRepo)
	tenantHandler := api.NewTenantHandler(tenantService)

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	authRouter := r.PathPrefix("/api").Subrouter()
	authRouter.Use(api.AuthMiddleware)

	authRouter.Handle("/tenants", api.GlobalAdminRequiredMiddleware(http.HandlerFunc(tenantHandler.CreateTenant))).Methods("POST")
	authRouter.Handle("/tenants", api.GlobalAdminRequiredMiddleware(http.HandlerFunc(tenantHandler.ListTenants))).Methods("GET")
	authRouter.Handle("/users/tenant-super-admin", api.GlobalAdminRequiredMiddleware(http.HandlerFunc(authHandler.CreateTenantSuperAdmin))).Methods("POST")
	
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	corsHandler := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)

	fmt.Println("Backend server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
