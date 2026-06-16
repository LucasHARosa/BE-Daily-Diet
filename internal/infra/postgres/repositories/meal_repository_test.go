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

func TestMealRepository_CreateAndGet(t *testing.T) {
	pool := testutils.NewTestPool(t)
	q := sqlc.New(pool)
	mealRepo := repositories.NewMealRepository(q)
	userRepo := repositories.NewUserRepository(q)
	ctx := context.Background()

	userID := uuid.New()
	email := "meal_test_" + userID.String() + "@example.com"
	if _, err := userRepo.Create(ctx, userID, "Test", email, "hash"); err != nil {
		t.Fatalf("setup: criar usuário falhou: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	desc := "Arroz e frango"
	calories := int32(500)
	meal, err := mealRepo.Create(ctx, uuid.New(), userID, "Almoço", &desc, time.Now(), true, &calories)
	if err != nil {
		t.Fatalf("Create() erro inesperado: %v", err)
	}
	if meal.Name != "Almoço" {
		t.Errorf("Create() Name = %v, esperado Almoço", meal.Name)
	}
	if meal.UserID != userID {
		t.Errorf("Create() UserID = %v, esperado %v", meal.UserID, userID)
	}
	if !meal.IsOnDiet {
		t.Error("Create() IsOnDiet deveria ser true")
	}

	found, err := mealRepo.GetByIDAndUserID(ctx, meal.ID, userID)
	if err != nil {
		t.Fatalf("GetByIDAndUserID() erro inesperado: %v", err)
	}
	if found.ID != meal.ID {
		t.Errorf("GetByIDAndUserID() ID = %v, esperado %v", found.ID, meal.ID)
	}
}

func TestMealRepository_List_WithFilters(t *testing.T) {
	pool := testutils.NewTestPool(t)
	q := sqlc.New(pool)
	mealRepo := repositories.NewMealRepository(q)
	userRepo := repositories.NewUserRepository(q)
	ctx := context.Background()

	userID := uuid.New()
	email := "meal_filter_" + userID.String() + "@example.com"
	if _, err := userRepo.Create(ctx, userID, "Filter Test", email, "hash"); err != nil {
		t.Fatalf("setup: criar usuário falhou: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	now := time.Now()
	if _, err := mealRepo.Create(ctx, uuid.New(), userID, "Refeição 1", nil, now, true, nil); err != nil {
		t.Fatalf("Create() refeição 1 falhou: %v", err)
	}
	if _, err := mealRepo.Create(ctx, uuid.New(), userID, "Refeição 2", nil, now, false, nil); err != nil {
		t.Fatalf("Create() refeição 2 falhou: %v", err)
	}

	all, err := mealRepo.List(ctx, userID, repositories.MealFilters{})
	if err != nil {
		t.Fatalf("List() sem filtros, erro inesperado: %v", err)
	}
	if len(all) < 2 {
		t.Errorf("List() retornou %d refeições, esperava ao menos 2", len(all))
	}

	onDiet := true
	filtered, err := mealRepo.List(ctx, userID, repositories.MealFilters{IsOnDiet: &onDiet})
	if err != nil {
		t.Fatalf("List() com filtro, erro inesperado: %v", err)
	}
	for _, m := range filtered {
		if !m.IsOnDiet {
			t.Errorf("List() com filtro on_diet retornou refeição fora da dieta: %v", m.Name)
		}
	}
}

func TestMealRepository_Update(t *testing.T) {
	pool := testutils.NewTestPool(t)
	q := sqlc.New(pool)
	mealRepo := repositories.NewMealRepository(q)
	userRepo := repositories.NewUserRepository(q)
	ctx := context.Background()

	userID := uuid.New()
	email := "meal_update_" + userID.String() + "@example.com"
	if _, err := userRepo.Create(ctx, userID, "Update Test", email, "hash"); err != nil {
		t.Fatalf("setup: criar usuário falhou: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	meal, err := mealRepo.Create(ctx, uuid.New(), userID, "Antes", nil, time.Now(), false, nil)
	if err != nil {
		t.Fatalf("Create() falhou: %v", err)
	}

	updated, err := mealRepo.Update(ctx, meal.ID, userID, "Depois", nil, time.Now(), true, nil)
	if err != nil {
		t.Fatalf("Update() erro inesperado: %v", err)
	}
	if updated.Name != "Depois" {
		t.Errorf("Update() Name = %v, esperado Depois", updated.Name)
	}
	if !updated.IsOnDiet {
		t.Error("Update() IsOnDiet deveria ser true após update")
	}
}

func TestMealRepository_Delete(t *testing.T) {
	pool := testutils.NewTestPool(t)
	q := sqlc.New(pool)
	mealRepo := repositories.NewMealRepository(q)
	userRepo := repositories.NewUserRepository(q)
	ctx := context.Background()

	userID := uuid.New()
	email := "meal_delete_" + userID.String() + "@example.com"
	if _, err := userRepo.Create(ctx, userID, "Delete Test", email, "hash"); err != nil {
		t.Fatalf("setup: criar usuário falhou: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	meal, err := mealRepo.Create(ctx, uuid.New(), userID, "Para deletar", nil, time.Now(), true, nil)
	if err != nil {
		t.Fatalf("Create() falhou: %v", err)
	}

	if err := mealRepo.Delete(ctx, meal.ID, userID); err != nil {
		t.Fatalf("Delete() erro inesperado: %v", err)
	}

	found, err := mealRepo.GetByIDAndUserID(ctx, meal.ID, userID)
	if err != nil {
		t.Fatalf("GetByIDAndUserID() após delete, erro inesperado: %v", err)
	}
	if found.ID != uuid.Nil {
		t.Errorf("GetByIDAndUserID() após delete deveria retornar vazio, retornou %v", found.ID)
	}
}

func TestMealRepository_GetByIDAndUserID_WrongUser(t *testing.T) {
	pool := testutils.NewTestPool(t)
	q := sqlc.New(pool)
	mealRepo := repositories.NewMealRepository(q)
	userRepo := repositories.NewUserRepository(q)
	ctx := context.Background()

	userID := uuid.New()
	email := "meal_wrong_" + userID.String() + "@example.com"
	if _, err := userRepo.Create(ctx, userID, "Wrong User Test", email, "hash"); err != nil {
		t.Fatalf("setup: criar usuário falhou: %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	})

	meal, err := mealRepo.Create(ctx, uuid.New(), userID, "Refeição", nil, time.Now(), true, nil)
	if err != nil {
		t.Fatalf("Create() falhou: %v", err)
	}

	// Tenta buscar a refeição com um user_id diferente
	otherUserID := uuid.New()
	found, err := mealRepo.GetByIDAndUserID(ctx, meal.ID, otherUserID)
	if err != nil {
		t.Fatalf("GetByIDAndUserID() erro inesperado: %v", err)
	}
	if found.ID != uuid.Nil {
		t.Error("GetByIDAndUserID() não deveria retornar refeição de outro usuário")
	}
}
