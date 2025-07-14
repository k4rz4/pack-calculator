package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"pack-calculator/internal/api/dto"
	apihttp "pack-calculator/internal/api/http"
	"pack-calculator/internal/domain/service"
	"pack-calculator/internal/infrastructure/logger"
)

// CalculationHandler handles calculation-related HTTP requests
type CalculationHandler struct {
	packService *service.PackService
	validator   *validator.Validate
}

// NewCalculationHandler creates a new calculation handler
func NewCalculationHandler(
	packService *service.PackService,
) *CalculationHandler {
	return &CalculationHandler{
		packService: packService,
		validator:   validator.New(),
	}
}

// Calculate handles POST /api/v1/calculate
func (h *CalculationHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := r.Header.Get("X-Request-ID")

	var req dto.CalculationRequest

	// Parse JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid JSON request", map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
		})
		apihttp.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	logger.Debug("Calculation request received", map[string]interface{}{
		"request_id":     requestID,
		"pack_sizes":     req.PackSizes,
		"order_quantity": req.OrderQuantity,
	})

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		logger.Warn("Request validation failed", map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
		})
		apihttp.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	// Perform calculation
	result, err := h.packService.CalculateOptimal(r.Context(), req.PackSizes, req.OrderQuantity)
	if err != nil {
		logger.Error("Calculation failed", map[string]interface{}{
			"request_id":     requestID,
			"pack_sizes":     req.PackSizes,
			"order_quantity": req.OrderQuantity,
			"error":          err.Error(),
		})
		apihttp.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Log successful calculation
	duration := time.Since(start)
	logger.Info("Calculation completed", map[string]interface{}{
		"request_id":     requestID,
		"pack_sizes":     req.PackSizes,
		"order_quantity": req.OrderQuantity,
		"total_items":    result.TotalItems,
		"total_packs":    result.TotalPacks,
		"items_overage":  result.ItemsOverage,
		"duration_ms":    duration.Milliseconds(),
		"calculation_id": result.ID,
	})

	// Convert to response DTO and return
	response := dto.ToCalculationResponse(result)
	apihttp.WriteSuccessResponse(w, http.StatusOK, response)
}
