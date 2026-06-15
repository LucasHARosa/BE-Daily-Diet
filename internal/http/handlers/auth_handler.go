package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/application/auth"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/middleware"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/responses"
	apperrors "github.com/LucasHARosa/BE-Daily-Diet/internal/shared/errors"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Email == "" || len(req.Password) < 8 {
		responses.Error(w, http.StatusBadRequest, "name, email and password (min 8 chars) are required")
		return
	}

	user, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			responses.Error(w, http.StatusConflict, "email already in use")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" {
		responses.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	result, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrUnauthorized) {
			responses.Error(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		responses.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	result, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, apperrors.ErrUnauthorized) {
			responses.Error(w, http.StatusUnauthorized, "invalid or expired refresh token")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	responses.JSON(w, http.StatusOK, map[string]string{"user_id": userID})
}
