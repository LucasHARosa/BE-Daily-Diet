package main

import (
	"context"
	"fmt"
	"net/http"

	appauth "github.com/LucasHARosa/BE-Daily-Diet/internal/application/auth"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/config"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/handlers"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/routes"
	infraauth "github.com/LucasHARosa/BE-Daily-Diet/internal/infra/auth"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/database"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/repositories"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/sqlc"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	ctx := context.Background()
	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	defer pool.Close()

	// Infra
	queries := sqlc.New(pool)
	jwtService := infraauth.NewJWTService(cfg.JWTSecret, cfg.JWTAccessTokenExpiresMin)

	// Repositories
	userRepo := repositories.NewUserRepository(queries)

	// Services
	authService := appauth.NewService(userRepo, jwtService, cfg.JWTRefreshTokenExpiresDays)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Router
	router := routes.Setup(authHandler, jwtService)

	fmt.Printf("Server running on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		panic(err)
	}
}
