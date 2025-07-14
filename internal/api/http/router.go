package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router handles HTTP routing configuration
type Router struct {
	router *mux.Router
}

// NewRouter creates a new HTTP router
func NewRouter() *Router {
	router := mux.NewRouter()

	return &Router{
		router: router,
	}
}

// RegisterCalculationRoutes registers calculation-related routes
func (r *Router) RegisterCalculationRoutes(calculateHandler http.HandlerFunc) {
	api := r.router.PathPrefix("/api/v1").Subrouter()

	// Calculation routes
	api.HandleFunc("/calculate", calculateHandler).Methods("POST")
}

// RegisterHealthRoutes registers health check routes
func (r *Router) RegisterHealthRoutes(healthHandler, readyHandler http.HandlerFunc) {
	// Health routes (allow both GET and HEAD for Docker healthcheck)
	r.router.HandleFunc("/health", healthHandler).Methods("GET", "HEAD")
	r.router.HandleFunc("/ready", readyHandler).Methods("GET", "HEAD")
}

// RegisterStaticRoutes registers static file routes
func (r *Router) RegisterStaticRoutes(uiHandler, staticHandler http.HandlerFunc) {
	// Serve UI at root path
	r.router.HandleFunc("/", uiHandler).Methods("GET")
	r.router.HandleFunc("/ui", uiHandler).Methods("GET")

	// Serve static assets
	r.router.PathPrefix("/static/").HandlerFunc(staticHandler).Methods("GET")
}

// Handler returns the underlying HTTP handler
func (r *Router) Handler() http.Handler {
	return r.router
}
