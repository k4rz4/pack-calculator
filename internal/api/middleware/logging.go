package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"pack-calculator/internal/infrastructure/logger"
)

// ResponseWriter wraps http.ResponseWriter to capture status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging middleware logs HTTP requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := generateRequestID()

		// Add request ID to context
		r.Header.Set("X-Request-ID", requestID)
		w.Header().Set("X-Request-ID", requestID)

		// Wrap response writer
		wrapped := &ResponseWriter{
			ResponseWriter: w,
			statusCode:     200, // default status
		}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Log the request
		duration := time.Since(start)

		fields := map[string]interface{}{
			"user_agent":    r.UserAgent(),
			"remote_addr":   r.RemoteAddr,
			"content_type":  r.Header.Get("Content-Type"),
			"response_size": w.Header().Get("Content-Length"),
		}

		logger.HTTP(r.Method, r.URL.Path, requestID, wrapped.statusCode, duration, fields)
	})
}

// generateRequestID creates a unique request ID
func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
