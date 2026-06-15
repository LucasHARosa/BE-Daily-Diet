package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/application/meals"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/middleware"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/responses"
	apperrors "github.com/LucasHARosa/BE-Daily-Diet/internal/shared/errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MealHandler struct {
	service *meals.Service
}

func NewMealHandler(service *meals.Service) *MealHandler {
	return &MealHandler{service: service}
}

type mealRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	EatenAt     string  `json:"eaten_at"`
	IsOnDiet    bool    `json:"is_on_diet"`
	Calories    *int32  `json:"calories"`
}

func (h *MealHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req mealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.EatenAt == "" {
		responses.Error(w, http.StatusBadRequest, "name and eaten_at are required")
		return
	}

	eatenAt, err := time.Parse(time.RFC3339, req.EatenAt)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "eaten_at must be in RFC3339 format (e.g. 2026-06-15T12:00:00Z)")
		return
	}

	meal, err := h.service.Create(r.Context(), userID, meals.CreateMealInput{
		Name:        req.Name,
		Description: req.Description,
		EatenAt:     eatenAt,
		IsOnDiet:    req.IsOnDiet,
		Calories:    req.Calories,
	})
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusCreated, meal)
}

func (h *MealHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	filters := meals.ListFilters{}

	if v := r.URL.Query().Get("start"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, "start must be in RFC3339 format")
			return
		}
		filters.StartDate = &t
	}
	if v := r.URL.Query().Get("end"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, "end must be in RFC3339 format")
			return
		}
		filters.EndDate = &t
	}
	if v := r.URL.Query().Get("status"); v != "" {
		switch v {
		case "on_diet":
			b := true
			filters.IsOnDiet = &b
		case "off_diet":
			b := false
			filters.IsOnDiet = &b
		default:
			responses.Error(w, http.StatusBadRequest, "status must be 'on_diet' or 'off_diet'")
			return
		}
	}

	list, err := h.service.List(r.Context(), userID, filters)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, list)
}

func (h *MealHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	mealID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid meal id")
		return
	}

	meal, err := h.service.GetByID(r.Context(), mealID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "meal not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, meal)
}

func (h *MealHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	mealID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid meal id")
		return
	}

	var req mealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.EatenAt == "" {
		responses.Error(w, http.StatusBadRequest, "name and eaten_at are required")
		return
	}

	eatenAt, err := time.Parse(time.RFC3339, req.EatenAt)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "eaten_at must be in RFC3339 format")
		return
	}

	meal, err := h.service.Update(r.Context(), mealID, userID, meals.UpdateMealInput{
		Name:        req.Name,
		Description: req.Description,
		EatenAt:     eatenAt,
		IsOnDiet:    req.IsOnDiet,
		Calories:    req.Calories,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "meal not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, meal)
}

func (h *MealHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	mealID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid meal id")
		return
	}

	if err := h.service.Delete(r.Context(), mealID, userID); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "meal not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getUserUUID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(middleware.GetUserID(r))
}
