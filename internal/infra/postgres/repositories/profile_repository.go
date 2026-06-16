package repositories

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserProfile struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	WeightKg            *float64
	HeightCm            *int32
	BirthDate           *string
	BodyFatPercentage   *float64
	BasalCalories       *int32
	ActivityLevel       *string
	GymFrequencyPerWeek *int32
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type ProfileRepository struct {
	q *sqlc.Queries
}

func NewProfileRepository(q *sqlc.Queries) *ProfileRepository {
	return &ProfileRepository{q: q}
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (UserProfile, error) {
	row, err := r.q.GetUserProfile(ctx, toPgUUID(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserProfile{}, nil
		}
		return UserProfile{}, err
	}
	return toUserProfile(row), nil
}

func (r *ProfileRepository) Upsert(ctx context.Context, userID uuid.UUID, p UserProfile) (UserProfile, error) {
	row, err := r.q.UpsertUserProfile(ctx, sqlc.UpsertUserProfileParams{
		ID:                  toPgUUID(uuid.New()),
		UserID:              toPgUUID(userID),
		WeightKg:            toPgNumeric(p.WeightKg),
		HeightCm:            toPgInt4(p.HeightCm),
		BirthDate:           toPgDate(p.BirthDate),
		BodyFatPercentage:   toPgNumeric(p.BodyFatPercentage),
		BasalCalories:       toPgInt4(p.BasalCalories),
		ActivityLevel:       toPgText(p.ActivityLevel),
		GymFrequencyPerWeek: toPgInt4(p.GymFrequencyPerWeek),
	})
	if err != nil {
		return UserProfile{}, err
	}
	return toUserProfile(row), nil
}

func toUserProfile(row sqlc.UserProfile) UserProfile {
	p := UserProfile{
		ID:        uuid.UUID(row.ID.Bytes),
		UserID:    uuid.UUID(row.UserID.Bytes),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if row.HeightCm.Valid {
		v := row.HeightCm.Int32
		p.HeightCm = &v
	}
	if row.BasalCalories.Valid {
		v := row.BasalCalories.Int32
		p.BasalCalories = &v
	}
	if row.GymFrequencyPerWeek.Valid {
		v := row.GymFrequencyPerWeek.Int32
		p.GymFrequencyPerWeek = &v
	}
	if row.ActivityLevel.Valid {
		p.ActivityLevel = &row.ActivityLevel.String
	}
	if row.BirthDate.Valid {
		s := row.BirthDate.Time.Format("2006-01-02")
		p.BirthDate = &s
	}
	p.WeightKg = fromPgNumeric(row.WeightKg)
	p.BodyFatPercentage = fromPgNumeric(row.BodyFatPercentage)
	return p
}

func toPgNumeric(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{}
	}
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(*f, 'f', 2, 64))
	return n
}

func fromPgNumeric(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return nil
	}
	v := f.Float64
	return &v
}

func toPgDate(s *string) pgtype.Date {
	if s == nil {
		return pgtype.Date{}
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: t, Valid: true}
}
