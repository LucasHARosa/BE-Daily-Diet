package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Meal struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description *string
	EatenAt     time.Time
	IsOnDiet    bool
	Calories    *int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MealFilters struct {
	StartDate *time.Time
	EndDate   *time.Time
	IsOnDiet  *bool
}

type MealRepository struct {
	q *sqlc.Queries
}

func NewMealRepository(q *sqlc.Queries) *MealRepository {
	return &MealRepository{q: q}
}

func (r *MealRepository) Create(ctx context.Context, id, userID uuid.UUID, name string, description *string, eatenAt time.Time, isOnDiet bool, calories *int32) (Meal, error) {
	row, err := r.q.CreateMeal(ctx, sqlc.CreateMealParams{
		ID:          toPgUUID(id),
		UserID:      toPgUUID(userID),
		Name:        name,
		Description: toPgText(description),
		EatenAt:     pgtype.Timestamp{Time: eatenAt, Valid: true},
		IsOnDiet:    isOnDiet,
		Calories:    toPgInt4(calories),
	})
	if err != nil {
		return Meal{}, err
	}
	return toMeal(row), nil
}

func (r *MealRepository) GetByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (Meal, error) {
	row, err := r.q.GetMealByIDAndUserID(ctx, sqlc.GetMealByIDAndUserIDParams{
		ID:     toPgUUID(id),
		UserID: toPgUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Meal{}, nil
		}
		return Meal{}, err
	}
	return toMeal(row), nil
}

func (r *MealRepository) List(ctx context.Context, userID uuid.UUID, filters MealFilters) ([]Meal, error) {
	params := sqlc.ListMealsByUserFilteredParams{
		UserID: toPgUUID(userID),
	}
	if filters.StartDate != nil {
		params.StartDate = pgtype.Timestamp{Time: *filters.StartDate, Valid: true}
	}
	if filters.EndDate != nil {
		params.EndDate = pgtype.Timestamp{Time: *filters.EndDate, Valid: true}
	}
	if filters.IsOnDiet != nil {
		params.IsOnDiet = pgtype.Bool{Bool: *filters.IsOnDiet, Valid: true}
	}

	rows, err := r.q.ListMealsByUserFiltered(ctx, params)
	if err != nil {
		return nil, err
	}

	meals := make([]Meal, len(rows))
	for i, row := range rows {
		meals[i] = toMeal(row)
	}
	return meals, nil
}

func (r *MealRepository) Update(ctx context.Context, id, userID uuid.UUID, name string, description *string, eatenAt time.Time, isOnDiet bool, calories *int32) (Meal, error) {
	row, err := r.q.UpdateMeal(ctx, sqlc.UpdateMealParams{
		ID:          toPgUUID(id),
		UserID:      toPgUUID(userID),
		Name:        name,
		Description: toPgText(description),
		EatenAt:     pgtype.Timestamp{Time: eatenAt, Valid: true},
		IsOnDiet:    isOnDiet,
		Calories:    toPgInt4(calories),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Meal{}, nil
		}
		return Meal{}, err
	}
	return toMeal(row), nil
}

func (r *MealRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.q.DeleteMeal(ctx, sqlc.DeleteMealParams{
		ID:     toPgUUID(id),
		UserID: toPgUUID(userID),
	})
}

func toMeal(row sqlc.Meal) Meal {
	m := Meal{
		ID:        uuid.UUID(row.ID.Bytes),
		UserID:    uuid.UUID(row.UserID.Bytes),
		Name:      row.Name,
		EatenAt:   row.EatenAt.Time,
		IsOnDiet:  row.IsOnDiet,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if row.Description.Valid {
		m.Description = &row.Description.String
	}
	if row.Calories.Valid {
		v := row.Calories.Int32
		m.Calories = &v
	}
	return m
}

func toPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func toPgInt4(n *int32) pgtype.Int4 {
	if n == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *n, Valid: true}
}
