package auth

import (
	"context"
	"fmt"
	"time"

	infraauth "github.com/LucasHARosa/BE-Daily-Diet/internal/infra/auth"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/repositories"
	apperrors "github.com/LucasHARosa/BE-Daily-Diet/internal/shared/errors"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type Service struct {
	repo       *repositories.UserRepository
	jwtService *infraauth.JWTService
	refreshExp int
}

func NewService(repo *repositories.UserRepository, jwtService *infraauth.JWTService, refreshExpDays int) *Service {
	return &Service{
		repo:       repo,
		jwtService: jwtService,
		refreshExp: refreshExpDays,
	}
}

func (s *Service) Register(ctx context.Context, name, email, password string) (UserResponse, error) {
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return UserResponse{}, err
	}
	if existing.ID != uuid.Nil {
		return UserResponse{}, fmt.Errorf("%w: email already in use", apperrors.ErrConflict)
	}

	hash, err := infraauth.HashPassword(password)
	if err != nil {
		return UserResponse{}, err
	}

	user, err := s.repo.Create(ctx, uuid.New(), name, email, hash)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return LoginResponse{}, err
	}
	if user.ID == uuid.Nil {
		return LoginResponse{}, apperrors.ErrUnauthorized
	}
	if !infraauth.CheckPassword(password, user.PasswordHash) {
		return LoginResponse{}, apperrors.ErrUnauthorized
	}

	accessToken, err := s.jwtService.Generate(user.ID)
	if err != nil {
		return LoginResponse{}, err
	}

	refreshToken := uuid.New().String()
	expiresAt := time.Now().Add(time.Duration(s.refreshExp) * 24 * time.Hour)
	err = s.repo.CreateRefreshToken(ctx, uuid.New(), user.ID, refreshToken, expiresAt)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, token string) (LoginResponse, error) {
	rt, err := s.repo.GetRefreshToken(ctx, token)
	if err != nil {
		return LoginResponse{}, err
	}
	if rt.Token == "" {
		return LoginResponse{}, apperrors.ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, rt.UserID)
	if err != nil || user.ID == uuid.Nil {
		return LoginResponse{}, apperrors.ErrUnauthorized
	}

	if err := s.repo.DeleteRefreshToken(ctx, token); err != nil {
		return LoginResponse{}, err
	}

	accessToken, err := s.jwtService.Generate(user.ID)
	if err != nil {
		return LoginResponse{}, err
	}

	newRefreshToken := uuid.New().String()
	expiresAt := time.Now().Add(time.Duration(s.refreshExp) * 24 * time.Hour)
	if err := s.repo.CreateRefreshToken(ctx, uuid.New(), user.ID, newRefreshToken, expiresAt); err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}
