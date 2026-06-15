package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/config"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/database"
	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Printf("Server running on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		panic(err)
	}
}
