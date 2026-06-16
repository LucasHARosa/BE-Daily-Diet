package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/repositories"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/sqlc"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/testutils"
	"github.com/google/uuid"
)

func TestUserRepository_Create(t *testing.T) {
	pool := testutils.NewTestPool(t)
	repo := repositories.NewUserRepository(sqlc.New(pool))
	ctx := context.Background()

	id := uuid.New()
	email := "test_" + id.String() + "@example.com"
	t.Cleanup(func() {
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	})

	user, err := repo.Create(ctx, id, "Test User", email, "hashed_password")
	if err != nil {
		t.Fatalf("Create() erro inesperado: %v", err)
	}
	if user.ID != id {
		t.Errorf("Create() ID = %v, esperado %v", user.ID, id)
	}
	if user.Email != email {
		t.Errorf("Create() Email = %v, esperado %v", user.Email, email)
	}
}

func TestUserRepository_GetByEmail_Exists(t *testing.T) {
	pool := testutils.NewTestPool(t)
	repo := repositories.NewUserRepository(sqlc.New(pool))
	ctx := context.Background()

	id := uuid.New()
	email := "test_" + id.String() + "@example.com"
	t.Cleanup(func() {
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	})

	if _, err := repo.Create(ctx, id, "Test User", email, "hashed_password"); err != nil {
		t.Fatalf("Create() setup falhou: %v", err)
	}

	found, err := repo.GetByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetByEmail() erro inesperado: %v", err)
	}
	if found.ID != id {
		t.Errorf("GetByEmail() ID = %v, esperado %v", found.ID, id)
	}
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	pool := testutils.NewTestPool(t)
	repo := repositories.NewUserRepository(sqlc.New(pool))
	ctx := context.Background()

	user, err := repo.GetByEmail(ctx, "nao_existe_"+uuid.NewString()+"@example.com")
	if err != nil {
		t.Fatalf("GetByEmail() erro inesperado: %v", err)
	}
	if user.ID != uuid.Nil {
		t.Errorf("GetByEmail() esperava ID nil para email inexistente, recebeu %v", user.ID)
	}
}

func TestUserRepository_RefreshToken(t *testing.T) {
	pool := testutils.NewTestPool(t)
	repo := repositories.NewUserRepository(sqlc.New(pool))
	ctx := context.Background()

	userID := uuid.New()
	email := "test_" + userID.String() + "@example.com"
	t.Cleanup(func() {
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	if _, err := repo.Create(ctx, userID, "Test User", email, "hash"); err != nil {
		t.Fatalf("Create() setup falhou: %v", err)
	}

	token := uuid.NewString()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := repo.CreateRefreshToken(ctx, uuid.New(), userID, token, expiresAt); err != nil {
		t.Fatalf("CreateRefreshToken() erro inesperado: %v", err)
	}

	rt, err := repo.GetRefreshToken(ctx, token)
	if err != nil {
		t.Fatalf("GetRefreshToken() erro inesperado: %v", err)
	}
	if rt.Token != token {
		t.Errorf("GetRefreshToken() Token = %v, esperado %v", rt.Token, token)
	}
	if rt.UserID != userID {
		t.Errorf("GetRefreshToken() UserID = %v, esperado %v", rt.UserID, userID)
	}

	if err := repo.DeleteRefreshToken(ctx, token); err != nil {
		t.Fatalf("DeleteRefreshToken() erro inesperado: %v", err)
	}

	deleted, err := repo.GetRefreshToken(ctx, token)
	if err != nil {
		t.Fatalf("GetRefreshToken() após delete, erro inesperado: %v", err)
	}
	if deleted.Token != "" {
		t.Errorf("GetRefreshToken() após delete deveria retornar vazio, retornou %v", deleted.Token)
	}
}
