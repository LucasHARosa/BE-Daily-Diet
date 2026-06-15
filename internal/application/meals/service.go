package meals

import (
	"context"
	"fmt"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/repositories"
	apperrors "github.com/LucasHARosa/BE-Daily-Diet/internal/shared/errors"
	"github.com/google/uuid"
)

type MealResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	EatenAt     time.Time `json:"eaten_at"`
	IsOnDiet    bool      `json:"is_on_diet"`
	Calories    *int32    `json:"calories"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateMealInput struct {
	Name        string
	Description *string
	EatenAt     time.Time
	IsOnDiet    bool
	Calories    *int32
}

type UpdateMealInput struct {
	Name        string
	Description *string
	EatenAt     time.Time
	IsOnDiet    bool
	Calories    *int32
}

type ListFilters struct {
	StartDate *time.Time
	EndDate   *time.Time
	IsOnDiet  *bool
}

type Service struct {
	repo *repositories.MealRepository
}

func NewService(repo *repositories.MealRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, input CreateMealInput) (MealResponse, error) {
	meal, err := s.repo.Create(ctx, uuid.New(), userID, input.Name, input.Description, input.EatenAt, input.IsOnDiet, input.Calories)
	if err != nil {
		return MealResponse{}, err
	}
	return toResponse(meal), nil
}

func (s *Service) List(ctx context.Context, userID uuid.UUID, filters ListFilters) ([]MealResponse, error) {
	meals, err := s.repo.List(ctx, userID, repositories.MealFilters{
		StartDate: filters.StartDate,
		EndDate:   filters.EndDate,
		IsOnDiet:  filters.IsOnDiet,
	})
	if err != nil {
		return nil, err
	}

	result := make([]MealResponse, len(meals))
	for i, m := range meals {
		result[i] = toResponse(m)
	}
	return result, nil
}

func (s *Service) GetByID(ctx context.Context, id, userID uuid.UUID) (MealResponse, error) {
	meal, err := s.repo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		return MealResponse{}, err
	}
	if meal.ID == uuid.Nil {
		return MealResponse{}, fmt.Errorf("%w: meal not found", apperrors.ErrNotFound)
	}
	return toResponse(meal), nil
}

func (s *Service) Update(ctx context.Context, id, userID uuid.UUID, input UpdateMealInput) (MealResponse, error) {
	existing, err := s.repo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		return MealResponse{}, err
	}
	if existing.ID == uuid.Nil {
		return MealResponse{}, fmt.Errorf("%w: meal not found", apperrors.ErrNotFound)
	}

	meal, err := s.repo.Update(ctx, id, userID, input.Name, input.Description, input.EatenAt, input.IsOnDiet, input.Calories)
	if err != nil {
		return MealResponse{}, err
	}
	return toResponse(meal), nil
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	existing, err := s.repo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	if existing.ID == uuid.Nil {
		return fmt.Errorf("%w: meal not found", apperrors.ErrNotFound)
	}
	return s.repo.Delete(ctx, id, userID)
}

func toResponse(m repositories.Meal) MealResponse {
	return MealResponse{
		ID:          m.ID,
		UserID:      m.UserID,
		Name:        m.Name,
		Description: m.Description,
		EatenAt:     m.EatenAt,
		IsOnDiet:    m.IsOnDiet,
		Calories:    m.Calories,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
