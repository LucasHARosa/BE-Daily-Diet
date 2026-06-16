package profile

import (
	"context"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/repositories"
	"github.com/google/uuid"
)

type ProfileOutput struct {
	WeightKg            *float64 `json:"weight_kg"`
	HeightCm            *int32   `json:"height_cm"`
	BirthDate           *string  `json:"birth_date"`
	BodyFatPercentage   *float64 `json:"body_fat_percentage"`
	BasalCalories       *int32   `json:"basal_calories"`
	ActivityLevel       *string  `json:"activity_level"`
	GymFrequencyPerWeek *int32   `json:"gym_frequency_per_week"`
}

type UpdateProfileInput struct {
	WeightKg            *float64
	HeightCm            *int32
	BirthDate           *string
	BodyFatPercentage   *float64
	BasalCalories       *int32
	ActivityLevel       *string
	GymFrequencyPerWeek *int32
}

type Service struct {
	repo *repositories.ProfileRepository
}

func NewService(repo *repositories.ProfileRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, userID uuid.UUID) (ProfileOutput, error) {
	p, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return ProfileOutput{}, err
	}
	return toOutput(p), nil
}

func (s *Service) Update(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (ProfileOutput, error) {
	p, err := s.repo.Upsert(ctx, userID, repositories.UserProfile{
		WeightKg:            input.WeightKg,
		HeightCm:            input.HeightCm,
		BirthDate:           input.BirthDate,
		BodyFatPercentage:   input.BodyFatPercentage,
		BasalCalories:       input.BasalCalories,
		ActivityLevel:       input.ActivityLevel,
		GymFrequencyPerWeek: input.GymFrequencyPerWeek,
	})
	if err != nil {
		return ProfileOutput{}, err
	}
	return toOutput(p), nil
}

func toOutput(p repositories.UserProfile) ProfileOutput {
	return ProfileOutput{
		WeightKg:            p.WeightKg,
		HeightCm:            p.HeightCm,
		BirthDate:           p.BirthDate,
		BodyFatPercentage:   p.BodyFatPercentage,
		BasalCalories:       p.BasalCalories,
		ActivityLevel:       p.ActivityLevel,
		GymFrequencyPerWeek: p.GymFrequencyPerWeek,
	}
}
