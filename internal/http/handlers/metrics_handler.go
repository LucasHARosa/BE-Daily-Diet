package handlers

import (
	"net/http"
	"time"

	"github.com/LucasHARosa/BE-Daily-Diet/internal/application/metrics"
	"github.com/LucasHARosa/BE-Daily-Diet/internal/http/responses"
)

type MetricsHandler struct {
	service *metrics.Service
}

func NewMetricsHandler(service *metrics.Service) *MetricsHandler {
	return &MetricsHandler{service: service}
}

func (h *MetricsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.service.GetSummary(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, result)
}

func (h *MetricsHandler) ByPeriod(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserUUID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	if startStr == "" || endStr == "" {
		responses.Error(w, http.StatusBadRequest, "start and end query params are required (e.g. 2026-06-01)")
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "start must be in YYYY-MM-DD format")
		return
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "end must be in YYYY-MM-DD format")
		return
	}
	// Inclui o dia inteiro do end
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	groupBy := r.URL.Query().Get("groupBy")
	if groupBy != "day" && groupBy != "month" {
		groupBy = "day"
	}

	result, err := h.service.GetByPeriod(r.Context(), userID, start, end, groupBy)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses.JSON(w, http.StatusOK, result)
}
