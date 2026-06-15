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

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) Create(ctx context.Context, id uuid.UUID, name, email, passwordHash string) (User, error) {
	row, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           toPgUUID(id),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return User{}, err
	}
	return toUser(row), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, nil
		}
		return User{}, err
	}
	return toUser(row), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	row, err := r.q.GetUserByID(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, nil
		}
		return User{}, err
	}
	return toUser(row), nil
}

func (r *UserRepository) CreateRefreshToken(ctx context.Context, tokenID, userID uuid.UUID, token string, expiresAt time.Time) error {
	_, err := r.q.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		ID:        toPgUUID(tokenID),
		UserID:    toPgUUID(userID),
		Token:     token,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	return err
}

func (r *UserRepository) GetRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row, err := r.q.GetRefreshToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RefreshToken{}, nil
		}
		return RefreshToken{}, err
	}
	return RefreshToken{
		ID:        uuid.UUID(row.ID.Bytes),
		UserID:    uuid.UUID(row.UserID.Bytes),
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt.Time,
	}, nil
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.q.DeleteRefreshToken(ctx, token)
}

func toUser(row sqlc.User) User {
	return User{
		ID:           uuid.UUID(row.ID.Bytes),
		Name:         row.Name,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}

func toPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}
