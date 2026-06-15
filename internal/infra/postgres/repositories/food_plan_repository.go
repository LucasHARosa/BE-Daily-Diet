package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FoodPlan struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FoodPlanDay struct {
	ID         uuid.UUID
	FoodPlanID uuid.UUID
	Weekday    int16
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type FoodPlanMeal struct {
	ID            uuid.UUID
	FoodPlanDayID uuid.UUID
	Name          string
	Description   *string
	ScheduledTime *string
	Calories      *int32
	SortOrder     int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type FoodPlanRepository struct {
	q *sqlc.Queries
}

func NewFoodPlanRepository(q *sqlc.Queries) *FoodPlanRepository {
	return &FoodPlanRepository{q: q}
}

// --- Food Plans ---

func (r *FoodPlanRepository) CreatePlan(ctx context.Context, id, userID uuid.UUID, title string, description *string) (FoodPlan, error) {
	row, err := r.q.CreateFoodPlan(ctx, sqlc.CreateFoodPlanParams{
		ID:          toPgUUID(id),
		UserID:      toPgUUID(userID),
		Title:       title,
		Description: toPgText(description),
		IsActive:    false,
	})
	if err != nil {
		return FoodPlan{}, err
	}
	return toFoodPlan(row), nil
}

func (r *FoodPlanRepository) GetPlanByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (FoodPlan, error) {
	row, err := r.q.GetFoodPlanByIDAndUserID(ctx, sqlc.GetFoodPlanByIDAndUserIDParams{
		ID:     toPgUUID(id),
		UserID: toPgUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FoodPlan{}, nil
		}
		return FoodPlan{}, err
	}
	return toFoodPlan(row), nil
}

func (r *FoodPlanRepository) ListPlansByUser(ctx context.Context, userID uuid.UUID) ([]FoodPlan, error) {
	rows, err := r.q.ListFoodPlansByUser(ctx, toPgUUID(userID))
	if err != nil {
		return nil, err
	}
	plans := make([]FoodPlan, len(rows))
	for i, row := range rows {
		plans[i] = toFoodPlan(row)
	}
	return plans, nil
}

func (r *FoodPlanRepository) GetActivePlan(ctx context.Context, userID uuid.UUID) (FoodPlan, error) {
	row, err := r.q.GetActiveFoodPlan(ctx, toPgUUID(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FoodPlan{}, nil
		}
		return FoodPlan{}, err
	}
	return toFoodPlan(row), nil
}

func (r *FoodPlanRepository) UpdatePlan(ctx context.Context, id, userID uuid.UUID, title string, description *string) (FoodPlan, error) {
	row, err := r.q.UpdateFoodPlan(ctx, sqlc.UpdateFoodPlanParams{
		ID:          toPgUUID(id),
		UserID:      toPgUUID(userID),
		Title:       title,
		Description: toPgText(description),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FoodPlan{}, nil
		}
		return FoodPlan{}, err
	}
	return toFoodPlan(row), nil
}

func (r *FoodPlanRepository) DeactivateAllUserPlans(ctx context.Context, userID uuid.UUID) error {
	return r.q.DeactivateAllUserFoodPlans(ctx, toPgUUID(userID))
}

func (r *FoodPlanRepository) SetPlanActive(ctx context.Context, id uuid.UUID, active bool) (FoodPlan, error) {
	row, err := r.q.SetFoodPlanActive(ctx, sqlc.SetFoodPlanActiveParams{
		ID:       toPgUUID(id),
		IsActive: active,
	})
	if err != nil {
		return FoodPlan{}, err
	}
	return toFoodPlan(row), nil
}

func (r *FoodPlanRepository) DeletePlan(ctx context.Context, id, userID uuid.UUID) error {
	return r.q.DeleteFoodPlan(ctx, sqlc.DeleteFoodPlanParams{
		ID:     toPgUUID(id),
		UserID: toPgUUID(userID),
	})
}

// --- Food Plan Days ---

func (r *FoodPlanRepository) CreateDay(ctx context.Context, id, foodPlanID uuid.UUID, weekday int16) (FoodPlanDay, error) {
	row, err := r.q.CreateFoodPlanDay(ctx, sqlc.CreateFoodPlanDayParams{
		ID:         toPgUUID(id),
		FoodPlanID: toPgUUID(foodPlanID),
		Weekday:    weekday,
	})
	if err != nil {
		return FoodPlanDay{}, err
	}
	return toFoodPlanDay(row), nil
}

func (r *FoodPlanRepository) GetDayByWeekday(ctx context.Context, foodPlanID uuid.UUID, weekday int16) (FoodPlanDay, error) {
	row, err := r.q.GetFoodPlanDayByWeekday(ctx, sqlc.GetFoodPlanDayByWeekdayParams{
		FoodPlanID: toPgUUID(foodPlanID),
		Weekday:    weekday,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FoodPlanDay{}, nil
		}
		return FoodPlanDay{}, err
	}
	return toFoodPlanDay(row), nil
}

func (r *FoodPlanRepository) ListDays(ctx context.Context, foodPlanID uuid.UUID) ([]FoodPlanDay, error) {
	rows, err := r.q.ListFoodPlanDays(ctx, toPgUUID(foodPlanID))
	if err != nil {
		return nil, err
	}
	days := make([]FoodPlanDay, len(rows))
	for i, row := range rows {
		days[i] = toFoodPlanDay(row)
	}
	return days, nil
}

// --- Food Plan Meals ---

func (r *FoodPlanRepository) CreateMealItem(ctx context.Context, id, dayID uuid.UUID, name string, description, scheduledTime *string, calories *int32, sortOrder int32) (FoodPlanMeal, error) {
	row, err := r.q.CreateFoodPlanMeal(ctx, sqlc.CreateFoodPlanMealParams{
		ID:            toPgUUID(id),
		FoodPlanDayID: toPgUUID(dayID),
		Name:          name,
		Description:   toPgText(description),
		ScheduledTime: toPgText(scheduledTime),
		Calories:      toPgInt4(calories),
		SortOrder:     sortOrder,
	})
	if err != nil {
		return FoodPlanMeal{}, err
	}
	return toFoodPlanMeal(row), nil
}

func (r *FoodPlanRepository) GetMealItemByID(ctx context.Context, id uuid.UUID) (FoodPlanMeal, error) {
	row, err := r.q.GetFoodPlanMealByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FoodPlanMeal{}, nil
		}
		return FoodPlanMeal{}, err
	}
	return toFoodPlanMeal(row), nil
}

func (r *FoodPlanRepository) ListMealItemsByDay(ctx context.Context, dayID uuid.UUID) ([]FoodPlanMeal, error) {
	rows, err := r.q.ListFoodPlanMealsByDay(ctx, toPgUUID(dayID))
	if err != nil {
		return nil, err
	}
	meals := make([]FoodPlanMeal, len(rows))
	for i, row := range rows {
		meals[i] = toFoodPlanMeal(row)
	}
	return meals, nil
}

func (r *FoodPlanRepository) UpdateMealItem(ctx context.Context, id uuid.UUID, name string, description, scheduledTime *string, calories *int32, sortOrder int32) (FoodPlanMeal, error) {
	row, err := r.q.UpdateFoodPlanMeal(ctx, sqlc.UpdateFoodPlanMealParams{
		ID:            toPgUUID(id),
		Name:          name,
		Description:   toPgText(description),
		ScheduledTime: toPgText(scheduledTime),
		Calories:      toPgInt4(calories),
		SortOrder:     sortOrder,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FoodPlanMeal{}, nil
		}
		return FoodPlanMeal{}, err
	}
	return toFoodPlanMeal(row), nil
}

func (r *FoodPlanRepository) DeleteMealItem(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteFoodPlanMeal(ctx, toPgUUID(id))
}

// --- converters ---

func toFoodPlan(row sqlc.FoodPlan) FoodPlan {
	fp := FoodPlan{
		ID:        uuid.UUID(row.ID.Bytes),
		UserID:    uuid.UUID(row.UserID.Bytes),
		Title:     row.Title,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if row.Description.Valid {
		fp.Description = &row.Description.String
	}
	return fp
}

func toFoodPlanDay(row sqlc.FoodPlanDay) FoodPlanDay {
	return FoodPlanDay{
		ID:         uuid.UUID(row.ID.Bytes),
		FoodPlanID: uuid.UUID(row.FoodPlanID.Bytes),
		Weekday:    row.Weekday,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}
}

func toFoodPlanMeal(row sqlc.FoodPlanMeal) FoodPlanMeal {
	m := FoodPlanMeal{
		ID:            uuid.UUID(row.ID.Bytes),
		FoodPlanDayID: uuid.UUID(row.FoodPlanDayID.Bytes),
		Name:          row.Name,
		SortOrder:     row.SortOrder,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
	if row.Description.Valid {
		m.Description = &row.Description.String
	}
	if row.ScheduledTime.Valid {
		m.ScheduledTime = &row.ScheduledTime.String
	}
	if row.Calories.Valid {
		v := row.Calories.Int32
		m.Calories = &v
	}
	return m
}

