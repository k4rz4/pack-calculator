package dto

// HealthResponse represents API response for health checks
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
	Time    string `json:"time"`
}
