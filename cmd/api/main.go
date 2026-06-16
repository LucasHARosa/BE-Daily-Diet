package main

import (
	"context"
	"fmt"
	"net/http"

	appauth "github.com/LucasHARosa/BE-Daily-Diet/internal/application/auth"
	appcalories "github.com/LucasHARosa/BE-Daily-Diet/internal/application/calories"
	appfoodplans "github.com/LucasHARosa/BE-Daily-Diet/internal/application/foodplans"
	appmeals "github.com/LucasHARosa/BE-Daily-Diet/internal/application/meals"
	appmetrics "github.com/LucasHARosa/BE-Daily-Diet/internal/application/metrics"
	appprofile "github.com/LucasHARosa/BE-Daily-Diet/internal/application/profile"
	appai "github.com/LucasHARosa/BE-Daily-Diet/internal/infra/ai"
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
	mealRepo := repositories.NewMealRepository(queries)
	metricsRepo := repositories.NewMetricsRepository(queries)
	foodPlanRepo := repositories.NewFoodPlanRepository(queries)
	profileRepo := repositories.NewProfileRepository(queries)
	calorieEstimationRepo := repositories.NewCalorieEstimationRepository(queries)

	// Infra
	calorieEstimator := appai.NewCalorieEstimator(cfg.AnthropicAPIKey)

	// Services
	authService := appauth.NewService(userRepo, jwtService, cfg.JWTRefreshTokenExpiresDays)
	mealService := appmeals.NewService(mealRepo)
	metricsService := appmetrics.NewService(metricsRepo)
	foodPlanService := appfoodplans.NewService(foodPlanRepo)
	profileService := appprofile.NewService(profileRepo)
	calorieService := appcalories.NewService(calorieEstimationRepo, calorieEstimator)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	mealHandler := handlers.NewMealHandler(mealService)
	metricsHandler := handlers.NewMetricsHandler(metricsService)
	foodPlanHandler := handlers.NewFoodPlanHandler(foodPlanService)
	profileHandler := handlers.NewProfileHandler(profileService)
	calorieHandler := handlers.NewCalorieEstimationHandler(calorieService)

	// Router
	router := routes.Setup(authHandler, mealHandler, metricsHandler, foodPlanHandler, profileHandler, calorieHandler, jwtService)

	fmt.Printf("Server running on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		panic(err)
	}
}
