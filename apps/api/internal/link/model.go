package link

import "time"

type Link struct {
	ID          string     `json:"id"`
	ShortCode   string     `json:"short_code"`
	Alias       string     `json:"alias,omitempty"`
	Destination string     `json:"destination"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type Status string

const (
	StatusActive   Status = "active"
	StatusDisabled Status = "disabled"
	StatusExpired  Status = "expired"
)
