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
	})

	return r
}
