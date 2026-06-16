package calories

import (
	"context"
	"time"

	appai "github.com/LucasHARosa/BE-Daily-Diet/internal/infra/ai"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/repositories"
	"github.com/google/uuid"
)

type EstimationItemOutput struct {
	Name     string `json:"name"`
	Calories int    `json:"calories"`
}

type EstimationOutput struct {
	ID                uuid.UUID              `json:"id"`
	Description       string                 `json:"description"`
	EstimatedCalories int32                  `json:"estimated_calories"`
	Confidence        string                 `json:"confidence"`
	Items             []EstimationItemOutput `json:"items"`
	Observation       *string                `json:"observation"`
	CreatedAt         time.Time              `json:"created_at"`
}

type Service struct {
	repo      *repositories.CalorieEstimationRepository
	estimator *appai.CalorieEstimator
}

func NewService(repo *repositories.CalorieEstimationRepository, estimator *appai.CalorieEstimator) *Service {
	return &Service{repo: repo, estimator: estimator}
}

func (s *Service) Estimate(ctx context.Context, userID uuid.UUID, description string) (EstimationOutput, error) {
	result, err := s.estimator.Estimate(ctx, description)
	if err != nil {
		return EstimationOutput{}, err
	}

	items := make([]repositories.EstimationItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = repositories.EstimationItem{Name: item.Name, Calories: item.Calories}
	}

	var observation *string
	if result.Observation != "" {
		observation = &result.Observation
	}

	est, err := s.repo.Create(ctx, userID, description, int32(result.EstimatedCalories), result.Confidence, items, observation)
	if err != nil {
		return EstimationOutput{}, err
	}

	return toOutput(est), nil
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]EstimationOutput, error) {
	estimations, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]EstimationOutput, len(estimations))
	for i, e := range estimations {
		out[i] = toOutput(e)
	}
	return out, nil
}

func toOutput(e repositories.CalorieEstimation) EstimationOutput {
	items := make([]EstimationItemOutput, len(e.Items))
	for i, item := range e.Items {
		items[i] = EstimationItemOutput{Name: item.Name, Calories: item.Calories}
	}
	return EstimationOutput{
		ID:                e.ID,
		Description:       e.Description,
		EstimatedCalories: e.EstimatedCalories,
		Confidence:        e.Confidence,
		Items:             items,
		Observation:       e.Observation,
		CreatedAt:         e.CreatedAt,
	}
}
