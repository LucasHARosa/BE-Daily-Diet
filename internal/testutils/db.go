package testutils

import (
	"context"
	"os"
	"testing"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	// Tenta carregar o .env subindo nos diretórios até encontrar
	paths := []string{".env", "../.env", "../../.env", "../../../.env", "../../../../.env"}
	for _, p := range paths {
		if err := godotenv.Load(p); err == nil {
			break
		}
	}
}

// NewTestPool cria uma conexão com o banco para testes de integração.
// Pula o teste automaticamente se DATABASE_URL não estiver configurada.
func NewTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL não configurada — pulando teste de integração")
	}
	pool, err := database.NewPostgresPool(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("falha ao conectar no banco de teste: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}
