package handlers

import (
	"net/http"
	"time"

	"pack-calculator/internal/api/dto"
	apihttp "pack-calculator/internal/api/http"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health handles GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := &dto.HealthResponse{
		Status:  "healthy",
		Version: "1.0.0", // TODO: Get from build info
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	apihttp.WriteJSONResponse(w, http.StatusOK, response)
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// TODO: Add database connectivity check
	// TODO: Add dependency health checks

	response := &dto.HealthResponse{
		Status: "ready",
		Time:   time.Now().UTC().Format(time.RFC3339),
	}

	apihttp.WriteJSONResponse(w, http.StatusOK, response)
}
