package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/application/calories"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/responses"
)

type CalorieEstimationHandler struct {
	service *calories.Service
}

func NewCalorieEstimationHandler(service *calories.Service) *CalorieEstimationHandler {
	return &CalorieEstimationHandler{service: service}
}

type estimateRequest struct {
	Description string `json:"description"`
}

func (h *CalorieEstimationHandler) Estimate(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req estimateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Description == "" {
		responses.Error(w, http.StatusBadRequest, "description is required")
		return
	}

	result, err := h.service.Estimate(r.Context(), userID, req.Description)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "failed to estimate calories")
		return
	}

	responses.JSON(w, http.StatusCreated, result)
}

func (h *CalorieEstimationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	estimations, err := h.service.ListByUser(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, estimations)
}
