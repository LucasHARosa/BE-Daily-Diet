package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type EstimationItem struct {
	Name     string `json:"name"`
	Calories int    `json:"calories"`
}

type CalorieEstimation struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	Description       string
	EstimatedCalories int32
	Confidence        string
	Items             []EstimationItem
	Observation       *string
	CreatedAt         time.Time
}

type CalorieEstimationRepository struct {
	q *sqlc.Queries
}

func NewCalorieEstimationRepository(q *sqlc.Queries) *CalorieEstimationRepository {
	return &CalorieEstimationRepository{q: q}
}

func (r *CalorieEstimationRepository) Create(ctx context.Context, userID uuid.UUID, description string, estimatedCalories int32, confidence string, items []EstimationItem, observation *string) (CalorieEstimation, error) {
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return CalorieEstimation{}, err
	}

	row, err := r.q.CreateCalorieEstimation(ctx, sqlc.CreateCalorieEstimationParams{
		ID:                toPgUUID(uuid.New()),
		UserID:            toPgUUID(userID),
		Description:       description,
		EstimatedCalories: estimatedCalories,
		Confidence:        confidence,
		Items:             itemsJSON,
		Observation:       toPgText(observation),
	})
	if err != nil {
		return CalorieEstimation{}, err
	}
	return toCalorieEstimation(row)
}

func (r *CalorieEstimationRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]CalorieEstimation, error) {
	rows, err := r.q.ListCalorieEstimationsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	result := make([]CalorieEstimation, 0, len(rows))
	for _, row := range rows {
		est, err := toCalorieEstimation(row)
		if err != nil {
			return nil, err
		}
		result = append(result, est)
	}
	return result, nil
}

func toCalorieEstimation(row sqlc.CalorieEstimation) (CalorieEstimation, error) {
	var items []EstimationItem
	if len(row.Items) > 0 {
		if err := json.Unmarshal(row.Items, &items); err != nil {
			return CalorieEstimation{}, err
		}
	}

	est := CalorieEstimation{
		ID:                uuid.UUID(row.ID.Bytes),
		UserID:            uuid.UUID(row.UserID.Bytes),
		Description:       row.Description,
		EstimatedCalories: row.EstimatedCalories,
		Confidence:        row.Confidence,
		Items:             items,
		CreatedAt:         row.CreatedAt.Time,
	}
	if row.Observation.Valid {
		est.Observation = &row.Observation.String
	}
	return est, nil
}
