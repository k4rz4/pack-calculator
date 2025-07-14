package http

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// WriteSuccessResponse writes a successful JSON response
func WriteSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: statusCode < 400,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// WriteErrorResponse writes an error JSON response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string, code ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorCode := ""
	if len(code) > 0 {
		errorCode = code[0]
	}

	response := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    errorCode,
	}

	json.NewEncoder(w).Encode(response)
}

// WriteJSONResponse writes a generic JSON response (for backward compatibility)
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
