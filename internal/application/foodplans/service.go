package foodplans

import (
	"context"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/repositories"
	apperrors "github.com/LucasHARosa/BE-Daily-Diet/internal/shared/errors"
	"github.com/google/uuid"
)

var weekdayNames = [7]string{
	"Segunda-feira", "Terça-feira", "Quarta-feira",
	"Quinta-feira", "Sexta-feira", "Sábado", "Domingo",
}

type MealItemOutput struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description"`
	ScheduledTime *string   `json:"scheduled_time"`
	Calories      *int32    `json:"calories"`
	SortOrder     int32     `json:"sort_order"`
}

type DayOutput struct {
	ID          uuid.UUID        `json:"id"`
	Weekday     int16            `json:"weekday"`
	WeekdayName string           `json:"weekday_name"`
	Meals       []MealItemOutput `json:"meals"`
}

type FoodPlanOutput struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
}

type FoodPlanDetailOutput struct {
	FoodPlanOutput
	Days []DayOutput `json:"days"`
}

type CreatePlanInput struct {
	Title       string
	Description *string
}

type UpdatePlanInput struct {
	Title       string
	Description *string
}

type CreateMealItemInput struct {
	Name          string
	Description   *string
	ScheduledTime *string
	Calories      *int32
	SortOrder     int32
}

type UpdateMealItemInput struct {
	Name          string
	Description   *string
	ScheduledTime *string
	Calories      *int32
	SortOrder     int32
}

type Service struct {
	repo *repositories.FoodPlanRepository
}

func NewService(repo *repositories.FoodPlanRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, input CreatePlanInput) (FoodPlanOutput, error) {
	plan, err := s.repo.CreatePlan(ctx, uuid.New(), userID, input.Title, input.Description)
	if err != nil {
		return FoodPlanOutput{}, err
	}
	return toFoodPlanOutput(plan), nil
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]FoodPlanOutput, error) {
	plans, err := s.repo.ListPlansByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]FoodPlanOutput, len(plans))
	for i, p := range plans {
		out[i] = toFoodPlanOutput(p)
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, id, userID uuid.UUID) (FoodPlanDetailOutput, error) {
	plan, err := s.repo.GetPlanByIDAndUserID(ctx, id, userID)
	if err != nil {
		return FoodPlanDetailOutput{}, err
	}
	if plan.ID == uuid.Nil {
		return FoodPlanDetailOutput{}, apperrors.ErrNotFound
	}
	return s.buildDetail(ctx, plan)
}

func (s *Service) GetActive(ctx context.Context, userID uuid.UUID) (FoodPlanDetailOutput, error) {
	plan, err := s.repo.GetActivePlan(ctx, userID)
	if err != nil {
		return FoodPlanDetailOutput{}, err
	}
	if plan.ID == uuid.Nil {
		return FoodPlanDetailOutput{}, apperrors.ErrNotFound
	}
	return s.buildDetail(ctx, plan)
}

func (s *Service) Update(ctx context.Context, id, userID uuid.UUID, input UpdatePlanInput) (FoodPlanOutput, error) {
	plan, err := s.repo.UpdatePlan(ctx, id, userID, input.Title, input.Description)
	if err != nil {
		return FoodPlanOutput{}, err
	}
	if plan.ID == uuid.Nil {
		return FoodPlanOutput{}, apperrors.ErrNotFound
	}
	return toFoodPlanOutput(plan), nil
}

func (s *Service) SetActive(ctx context.Context, id, userID uuid.UUID, active bool) (FoodPlanOutput, error) {
	existing, err := s.repo.GetPlanByIDAndUserID(ctx, id, userID)
	if err != nil {
		return FoodPlanOutput{}, err
	}
	if existing.ID == uuid.Nil {
		return FoodPlanOutput{}, apperrors.ErrNotFound
	}
	if active {
		if err := s.repo.DeactivateAllUserPlans(ctx, userID); err != nil {
			return FoodPlanOutput{}, err
		}
	}
	plan, err := s.repo.SetPlanActive(ctx, id, active)
	if err != nil {
		return FoodPlanOutput{}, err
	}
	return toFoodPlanOutput(plan), nil
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	existing, err := s.repo.GetPlanByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	if existing.ID == uuid.Nil {
		return apperrors.ErrNotFound
	}
	return s.repo.DeletePlan(ctx, id, userID)
}

func (s *Service) AddMealToDay(ctx context.Context, planID, userID uuid.UUID, weekday int16, input CreateMealItemInput) (MealItemOutput, error) {
	plan, err := s.repo.GetPlanByIDAndUserID(ctx, planID, userID)
	if err != nil {
		return MealItemOutput{}, err
	}
	if plan.ID == uuid.Nil {
		return MealItemOutput{}, apperrors.ErrNotFound
	}

	day, err := s.repo.GetDayByWeekday(ctx, planID, weekday)
	if err != nil {
		return MealItemOutput{}, err
	}
	if day.ID == uuid.Nil {
		day, err = s.repo.CreateDay(ctx, uuid.New(), planID, weekday)
		if err != nil {
			return MealItemOutput{}, err
		}
	}

	meal, err := s.repo.CreateMealItem(ctx, uuid.New(), day.ID, input.Name, input.Description, input.ScheduledTime, input.Calories, input.SortOrder)
	if err != nil {
		return MealItemOutput{}, err
	}
	return toMealItemOutput(meal), nil
}

func (s *Service) UpdateMealItem(ctx context.Context, mealID uuid.UUID, input UpdateMealItemInput) (MealItemOutput, error) {
	existing, err := s.repo.GetMealItemByID(ctx, mealID)
	if err != nil {
		return MealItemOutput{}, err
	}
	if existing.ID == uuid.Nil {
		return MealItemOutput{}, apperrors.ErrNotFound
	}
	updated, err := s.repo.UpdateMealItem(ctx, mealID, input.Name, input.Description, input.ScheduledTime, input.Calories, input.SortOrder)
	if err != nil {
		return MealItemOutput{}, err
	}
	return toMealItemOutput(updated), nil
}

func (s *Service) DeleteMealItem(ctx context.Context, mealID uuid.UUID) error {
	existing, err := s.repo.GetMealItemByID(ctx, mealID)
	if err != nil {
		return err
	}
	if existing.ID == uuid.Nil {
		return apperrors.ErrNotFound
	}
	return s.repo.DeleteMealItem(ctx, mealID)
}

func (s *Service) buildDetail(ctx context.Context, plan repositories.FoodPlan) (FoodPlanDetailOutput, error) {
	days, err := s.repo.ListDays(ctx, plan.ID)
	if err != nil {
		return FoodPlanDetailOutput{}, err
	}

	dayOutputs := make([]DayOutput, len(days))
	for i, day := range days {
		meals, err := s.repo.ListMealItemsByDay(ctx, day.ID)
		if err != nil {
			return FoodPlanDetailOutput{}, err
		}
		mealOutputs := make([]MealItemOutput, len(meals))
		for j, m := range meals {
			mealOutputs[j] = toMealItemOutput(m)
		}
		weekdayName := ""
		if int(day.Weekday) < len(weekdayNames) {
			weekdayName = weekdayNames[day.Weekday]
		}
		dayOutputs[i] = DayOutput{
			ID:          day.ID,
			Weekday:     day.Weekday,
			WeekdayName: weekdayName,
			Meals:       mealOutputs,
		}
	}

	return FoodPlanDetailOutput{
		FoodPlanOutput: toFoodPlanOutput(plan),
		Days:           dayOutputs,
	}, nil
}

func toFoodPlanOutput(p repositories.FoodPlan) FoodPlanOutput {
	return FoodPlanOutput{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		IsActive:    p.IsActive,
	}
}

func toMealItemOutput(m repositories.FoodPlanMeal) MealItemOutput {
	return MealItemOutput{
		ID:            m.ID,
		Name:          m.Name,
		Description:   m.Description,
		ScheduledTime: m.ScheduledTime,
		Calories:      m.Calories,
		SortOrder:     m.SortOrder,
	}
}
