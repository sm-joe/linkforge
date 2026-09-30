package link

import "time"

func (l *Link) IsExpired(now time.Time) bool {
	if l.ExpiresAt == nil {
		return false
	}

	return !now.Before(l.ExpiresAt.UTC())
}

func (l *Link) IsAvailable(now time.Time) bool {
	if l.Status != StatusActive {
		return false
	}

	return !l.IsExpired(now)
}
