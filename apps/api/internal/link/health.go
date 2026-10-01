package link

import "time"

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

type LinkHealth struct {
	ShortCode      string       `json:"short_code"`
	Destination    string       `json:"destination"`
	Status         HealthStatus `json:"status"`
	HTTPStatus     int          `json:"http_status,omitempty"`
	ResponseTimeMS int64        `json:"response_time_ms"`
	FinalURL       string       `json:"final_url,omitempty"`
	Redirects      int          `json:"redirects"`
	HTTPS          bool         `json:"https"`
	CheckedAt      time.Time    `json:"checked_at"`
	Error          string       `json:"error,omitempty"`
}
