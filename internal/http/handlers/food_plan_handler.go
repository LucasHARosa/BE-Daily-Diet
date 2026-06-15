package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/application/foodplans"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/responses"
	apperrors "github.com/LucasHARosa/BE-Daily-Diet/internal/shared/errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FoodPlanHandler struct {
	service *foodplans.Service
}

func NewFoodPlanHandler(service *foodplans.Service) *FoodPlanHandler {
	return &FoodPlanHandler{service: service}
}

type createPlanRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type updatePlanRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type setActiveRequest struct {
	IsActive bool `json:"is_active"`
}

type createMealItemRequest struct {
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	ScheduledTime *string `json:"scheduled_time"`
	Calories      *int32  `json:"calories"`
	SortOrder     int32   `json:"sort_order"`
}

func (h *FoodPlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		responses.Error(w, http.StatusBadRequest, "title is required")
		return
	}

	plan, err := h.service.Create(r.Context(), userID, foodplans.CreatePlanInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusCreated, plan)
}

func (h *FoodPlanHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plans, err := h.service.List(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, plans)
}

func (h *FoodPlanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	planID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid food plan id")
		return
	}

	plan, err := h.service.GetByID(r.Context(), planID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "food plan not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, plan)
}

func (h *FoodPlanHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plan, err := h.service.GetActive(r.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "no active food plan")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, plan)
}

func (h *FoodPlanHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	planID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid food plan id")
		return
	}

	var req updatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		responses.Error(w, http.StatusBadRequest, "title is required")
		return
	}

	plan, err := h.service.Update(r.Context(), planID, userID, foodplans.UpdatePlanInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "food plan not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, plan)
}

func (h *FoodPlanHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	planID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid food plan id")
		return
	}

	var req setActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	plan, err := h.service.SetActive(r.Context(), planID, userID, req.IsActive)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "food plan not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, plan)
}

func (h *FoodPlanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	planID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid food plan id")
		return
	}

	if err := h.service.Delete(r.Context(), planID, userID); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "food plan not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FoodPlanHandler) AddMealToDay(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	planID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid food plan id")
		return
	}

	weekdayStr := chi.URLParam(r, "weekday")
	weekdayInt, err := strconv.ParseInt(weekdayStr, 10, 16)
	if err != nil || weekdayInt < 0 || weekdayInt > 6 {
		responses.Error(w, http.StatusBadRequest, "weekday must be 0 (Monday) to 6 (Sunday)")
		return
	}

	var req createMealItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		responses.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	meal, err := h.service.AddMealToDay(r.Context(), planID, userID, int16(weekdayInt), foodplans.CreateMealItemInput{
		Name:          req.Name,
		Description:   req.Description,
		ScheduledTime: req.ScheduledTime,
		Calories:      req.Calories,
		SortOrder:     req.SortOrder,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "food plan not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusCreated, meal)
}

func (h *FoodPlanHandler) UpdateMealItem(w http.ResponseWriter, r *http.Request) {
	_, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	mealID, err := uuid.Parse(chi.URLParam(r, "mealId"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid meal id")
		return
	}

	var req createMealItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		responses.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	meal, err := h.service.UpdateMealItem(r.Context(), mealID, foodplans.UpdateMealItemInput{
		Name:          req.Name,
		Description:   req.Description,
		ScheduledTime: req.ScheduledTime,
		Calories:      req.Calories,
		SortOrder:     req.SortOrder,
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

func (h *FoodPlanHandler) DeleteMealItem(w http.ResponseWriter, r *http.Request) {
	_, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	mealID, err := uuid.Parse(chi.URLParam(r, "mealId"))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid meal id")
		return
	}

	if err := h.service.DeleteMealItem(r.Context(), mealID); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			responses.Error(w, http.StatusNotFound, "meal not found")
			return
		}
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
