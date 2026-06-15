package metrics

import (
	"context"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/domain"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/repositories"
	"github.com/google/uuid"
)

type SummaryResponse struct {
	TotalMeals         int64   `json:"totalMeals"`
	TotalOnDiet        int64   `json:"totalOnDiet"`
	TotalOffDiet       int64   `json:"totalOffDiet"`
	BestOnDietSequence int     `json:"bestOnDietSequence"`
	OnDietPercentage   float64 `json:"onDietPercentage"`
}

type PeriodSummary struct {
	TotalMeals         int64   `json:"totalMeals"`
	TotalOnDiet        int64   `json:"totalOnDiet"`
	TotalOffDiet       int64   `json:"totalOffDiet"`
	TotalCalories      int64   `json:"totalCalories"`
	BestOnDietSequence int     `json:"bestOnDietSequence"`
}

type GroupEntry struct {
	Date          string `json:"date"`
	TotalMeals    int64  `json:"totalMeals"`
	TotalOnDiet   int64  `json:"totalOnDiet"`
	TotalOffDiet  int64  `json:"totalOffDiet"`
	TotalCalories int64  `json:"totalCalories"`
}

type PeriodResponse struct {
	Period  PeriodRange   `json:"period"`
	Summary PeriodSummary `json:"summary"`
	Groups  []GroupEntry  `json:"groups"`
}

type PeriodRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Service struct {
	repo *repositories.MetricsRepository
}

func NewService(repo *repositories.MetricsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSummary(ctx context.Context, userID uuid.UUID) (SummaryResponse, error) {
	totals, err := s.repo.GetSummary(ctx, userID)
	if err != nil {
		return SummaryResponse{}, err
	}

	onDietList, err := s.repo.ListOnDietStatus(ctx, userID, nil, nil)
	if err != nil {
		return SummaryResponse{}, err
	}

	bestStreak := domain.CalculateBestStreak(onDietList)

	var pct float64
	if totals.TotalMeals > 0 {
		pct = float64(totals.TotalOnDiet) / float64(totals.TotalMeals) * 100
		pct = roundFloat(pct, 2)
	}

	return SummaryResponse{
		TotalMeals:         totals.TotalMeals,
		TotalOnDiet:        totals.TotalOnDiet,
		TotalOffDiet:       totals.TotalOffDiet,
		BestOnDietSequence: bestStreak,
		OnDietPercentage:   pct,
	}, nil
}

func (s *Service) GetByPeriod(ctx context.Context, userID uuid.UUID, start, end time.Time, groupBy string) (PeriodResponse, error) {
	var groups []repositories.MealGroupRow
	var err error

	if groupBy == "month" {
		groups, err = s.repo.ListGroupedByMonth(ctx, userID, start, end)
	} else {
		groups, err = s.repo.ListGroupedByDay(ctx, userID, start, end)
	}
	if err != nil {
		return PeriodResponse{}, err
	}

	onDietList, err := s.repo.ListOnDietStatus(ctx, userID, &start, &end)
	if err != nil {
		return PeriodResponse{}, err
	}

	bestStreak := domain.CalculateBestStreak(onDietList)

	var summary PeriodSummary
	summary.BestOnDietSequence = bestStreak

	entries := make([]GroupEntry, len(groups))
	for i, g := range groups {
		summary.TotalMeals += g.TotalMeals
		summary.TotalOnDiet += g.TotalOnDiet
		summary.TotalOffDiet += g.TotalOffDiet
		summary.TotalCalories += g.TotalCalories

		entries[i] = GroupEntry{
			Date:          g.Date.Format("2006-01-02"),
			TotalMeals:    g.TotalMeals,
			TotalOnDiet:   g.TotalOnDiet,
			TotalOffDiet:  g.TotalOffDiet,
			TotalCalories: g.TotalCalories,
		}
	}

	dateFormat := "2006-01-02"
	return PeriodResponse{
		Period: PeriodRange{
			Start: start.Format(dateFormat),
			End:   end.Format(dateFormat),
		},
		Summary: summary,
		Groups:  entries,
	}, nil
}

func roundFloat(val float64, precision int) float64 {
	ratio := 1.0
	for range precision {
		ratio *= 10
	}
	return float64(int(val*ratio+0.5)) / ratio
}
