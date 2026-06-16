package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/application/profile"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/responses"
)

type ProfileHandler struct {
	service *profile.Service
}

func NewProfileHandler(service *profile.Service) *ProfileHandler {
	return &ProfileHandler{service: service}
}

type updateProfileRequest struct {
	WeightKg            *float64 `json:"weight_kg"`
	HeightCm            *int32   `json:"height_cm"`
	BirthDate           *string  `json:"birth_date"`
	BodyFatPercentage   *float64 `json:"body_fat_percentage"`
	BasalCalories       *int32   `json:"basal_calories"`
	ActivityLevel       *string  `json:"activity_level"`
	GymFrequencyPerWeek *int32   `json:"gym_frequency_per_week"`
}

func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	p, err := h.service.Get(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, p)
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.service.Update(r.Context(), userID, profile.UpdateProfileInput{
		WeightKg:            req.WeightKg,
		HeightCm:            req.HeightCm,
		BirthDate:           req.BirthDate,
		BodyFatPercentage:   req.BodyFatPercentage,
		BasalCalories:       req.BasalCalories,
		ActivityLevel:       req.ActivityLevel,
		GymFrequencyPerWeek: req.GymFrequencyPerWeek,
	})
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, p)
}
