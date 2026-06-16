package auth_test

import (
	"context"
	"errors"
	"testing"

	appauth "github.com/LucasHARosa/BE-Daily-Diet/internal/application/auth"
	infraauth "github.com/LucasHARosa/BE-Daily-Diet/internal/infra/auth"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/repositories"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/infra/postgres/sqlc"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/testutils"
	apperrors "github.com/LucasHARosa/BE-Daily-Diet/internal/shared/errors"
	"github.com/google/uuid"
)

func newAuthService(t *testing.T) (*appauth.Service, func()) {
	t.Helper()
	pool := testutils.NewTestPool(t)
	q := sqlc.New(pool)
	userRepo := repositories.NewUserRepository(q)
	jwtService := infraauth.NewJWTService("test-secret", 30)
	svc := appauth.NewService(userRepo, jwtService, 7)
	cleanup := func() { pool.Close() }
	return svc, cleanup
}

func TestAuthService_Register_Success(t *testing.T) {
	svc, cleanup := newAuthService(t)
	defer cleanup()
	ctx := context.Background()

	email := "auth_test_" + uuid.NewString() + "@example.com"
	user, err := svc.Register(ctx, "Lucas", email, "senha123")
	if err != nil {
		t.Fatalf("Register() erro inesperado: %v", err)
	}
	if user.Email != email {
		t.Errorf("Register() Email = %v, esperado %v", user.Email, email)
	}
	if user.ID == uuid.Nil {
		t.Error("Register() ID não deveria ser nil")
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	svc, cleanup := newAuthService(t)
	defer cleanup()
	ctx := context.Background()

	email := "dup_" + uuid.NewString() + "@example.com"
	if _, err := svc.Register(ctx, "Lucas", email, "senha123"); err != nil {
		t.Fatalf("Register() primeira vez falhou: %v", err)
	}

	_, err := svc.Register(ctx, "Lucas", email, "senha123")
	if err == nil {
		t.Fatal("Register() deveria retornar erro para email duplicado")
	}
	if !errors.Is(err, apperrors.ErrConflict) {
		t.Errorf("Register() erro = %v, esperado ErrConflict", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, cleanup := newAuthService(t)
	defer cleanup()
	ctx := context.Background()

	email := "login_" + uuid.NewString() + "@example.com"
	if _, err := svc.Register(ctx, "Lucas", email, "senha123"); err != nil {
		t.Fatalf("Register() setup falhou: %v", err)
	}

	resp, err := svc.Login(ctx, email, "senha123")
	if err != nil {
		t.Fatalf("Login() erro inesperado: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("Login() AccessToken não deveria estar vazio")
	}
	if resp.RefreshToken == "" {
		t.Error("Login() RefreshToken não deveria estar vazio")
	}
	if resp.User.Email != email {
		t.Errorf("Login() User.Email = %v, esperado %v", resp.User.Email, email)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc, cleanup := newAuthService(t)
	defer cleanup()
	ctx := context.Background()

	email := "wrongpw_" + uuid.NewString() + "@example.com"
	if _, err := svc.Register(ctx, "Lucas", email, "senha123"); err != nil {
		t.Fatalf("Register() setup falhou: %v", err)
	}

	_, err := svc.Login(ctx, email, "senha_errada")
	if err == nil {
		t.Fatal("Login() deveria retornar erro para senha errada")
	}
	if !errors.Is(err, apperrors.ErrUnauthorized) {
		t.Errorf("Login() erro = %v, esperado ErrUnauthorized", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc, cleanup := newAuthService(t)
	defer cleanup()
	ctx := context.Background()

	_, err := svc.Login(ctx, "nao_existe@example.com", "senha")
	if err == nil {
		t.Fatal("Login() deveria retornar erro para usuário inexistente")
	}
	if !errors.Is(err, apperrors.ErrUnauthorized) {
		t.Errorf("Login() erro = %v, esperado ErrUnauthorized", err)
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	svc, cleanup := newAuthService(t)
	defer cleanup()
	ctx := context.Background()

	email := "refresh_" + uuid.NewString() + "@example.com"
	if _, err := svc.Register(ctx, "Lucas", email, "senha123"); err != nil {
		t.Fatalf("Register() setup falhou: %v", err)
	}

	login, err := svc.Login(ctx, email, "senha123")
	if err != nil {
		t.Fatalf("Login() setup falhou: %v", err)
	}

	refreshed, err := svc.RefreshToken(ctx, login.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshToken() erro inesperado: %v", err)
	}
	if refreshed.AccessToken == "" {
		t.Error("RefreshToken() novo AccessToken não deveria estar vazio")
	}
	// O refresh token antigo deve ser invalidado (rotação de tokens)
	_, err = svc.RefreshToken(ctx, login.RefreshToken)
	if err == nil {
		t.Error("RefreshToken() deveria falhar ao reutilizar o mesmo refresh token")
	}
}
