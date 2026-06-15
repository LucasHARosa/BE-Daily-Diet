package routes

import (
	"net/http"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/handlers"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/middleware"
	infraauth "github.com/LucasHARosa/BE-Daily-Diet/internal/infra/auth"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func Setup(
	authHandler *handlers.AuthHandler,
	mealHandler *handlers.MealHandler,
	metricsHandler *handlers.MetricsHandler,
	foodPlanHandler *handlers.FoodPlanHandler,
	jwtService *infraauth.JWTService,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Rotas públicas
	r.Post("/users", authHandler.Register)
	r.Post("/sessions", authHandler.Login)
	r.Post("/sessions/refresh", authHandler.RefreshToken)

	// Rotas protegidas
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtService))

		r.Get("/me", authHandler.Me)

		r.Post("/meals", mealHandler.Create)
		r.Get("/meals", mealHandler.List)
		r.Get("/meals/{id}", mealHandler.GetByID)
		r.Put("/meals/{id}", mealHandler.Update)
		r.Delete("/meals/{id}", mealHandler.Delete)

		r.Get("/metrics/summary", metricsHandler.Summary)
		r.Get("/metrics", metricsHandler.ByPeriod)

		r.Post("/food-plans", foodPlanHandler.Create)
		r.Get("/food-plans", foodPlanHandler.List)
		r.Get("/food-plans/active", foodPlanHandler.GetActive)
		r.Get("/food-plans/{id}", foodPlanHandler.GetByID)
		r.Put("/food-plans/{id}", foodPlanHandler.Update)
		r.Patch("/food-plans/{id}/active", foodPlanHandler.SetActive)
		r.Delete("/food-plans/{id}", foodPlanHandler.Delete)

		r.Post("/food-plans/{id}/days/{weekday}/meals", foodPlanHandler.AddMealToDay)
		r.Put("/food-plan-meals/{mealId}", foodPlanHandler.UpdateMealItem)
		r.Delete("/food-plan-meals/{mealId}", foodPlanHandler.DeleteMealItem)
	})

	return r
}
