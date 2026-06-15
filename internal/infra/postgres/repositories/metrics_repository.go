package repositories

import (
	"context"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type MetricsSummaryRow struct {
	TotalMeals   int64
	TotalOnDiet  int64
	TotalOffDiet int64
}

type MealGroupRow struct {
	Date          time.Time
	TotalMeals    int64
	TotalOnDiet   int64
	TotalOffDiet  int64
	TotalCalories int64
}

type MetricsRepository struct {
	q *sqlc.Queries
}

func NewMetricsRepository(q *sqlc.Queries) *MetricsRepository {
	return &MetricsRepository{q: q}
}

func (r *MetricsRepository) GetSummary(ctx context.Context, userID uuid.UUID) (MetricsSummaryRow, error) {
	row, err := r.q.GetMetricsSummary(ctx, toPgUUID(userID))
	if err != nil {
		return MetricsSummaryRow{}, err
	}
	return MetricsSummaryRow{
		TotalMeals:   row.TotalMeals,
		TotalOnDiet:  row.TotalOnDiet,
		TotalOffDiet: row.TotalOffDiet,
	}, nil
}

func (r *MetricsRepository) ListOnDietStatus(ctx context.Context, userID uuid.UUID, start, end *time.Time) ([]bool, error) {
	params := sqlc.ListMealsOnDietStatusParams{
		UserID: toPgUUID(userID),
	}
	if start != nil {
		params.StartDate = pgtype.Timestamp{Time: *start, Valid: true}
	}
	if end != nil {
		params.EndDate = pgtype.Timestamp{Time: *end, Valid: true}
	}
	return r.q.ListMealsOnDietStatus(ctx, params)
}

func (r *MetricsRepository) ListGroupedByDay(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]MealGroupRow, error) {
	rows, err := r.q.ListMealsGroupedByDay(ctx, sqlc.ListMealsGroupedByDayParams{
		UserID:    toPgUUID(userID),
		EatenAt:   pgtype.Timestamp{Time: start, Valid: true},
		EatenAt_2: pgtype.Timestamp{Time: end, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return dayRowsToGroup(rows), nil
}

func (r *MetricsRepository) ListGroupedByMonth(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]MealGroupRow, error) {
	rows, err := r.q.ListMealsGroupedByMonth(ctx, sqlc.ListMealsGroupedByMonthParams{
		UserID:    toPgUUID(userID),
		EatenAt:   pgtype.Timestamp{Time: start, Valid: true},
		EatenAt_2: pgtype.Timestamp{Time: end, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return monthRowsToGroup(rows), nil
}

func dayRowsToGroup(rows []sqlc.ListMealsGroupedByDayRow) []MealGroupRow {
	result := make([]MealGroupRow, len(rows))
	for i, row := range rows {
		result[i] = MealGroupRow{
			Date:          row.Date.Time,
			TotalMeals:    row.TotalMeals,
			TotalOnDiet:   row.TotalOnDiet,
			TotalOffDiet:  row.TotalOffDiet,
			TotalCalories: row.TotalCalories,
		}
	}
	return result
}

func monthRowsToGroup(rows []sqlc.ListMealsGroupedByMonthRow) []MealGroupRow {
	result := make([]MealGroupRow, len(rows))
	for i, row := range rows {
		result[i] = MealGroupRow{
			Date:          row.Date.Time,
			TotalMeals:    row.TotalMeals,
			TotalOnDiet:   row.TotalOnDiet,
			TotalOffDiet:  row.TotalOffDiet,
			TotalCalories: row.TotalCalories,
		}
	}
	return result
}
