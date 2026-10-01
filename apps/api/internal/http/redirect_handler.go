package httpserver

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sm-joe/linkforge/internal/link"
)

type RedirectHandler struct {
	service *link.Service
	clicks  *link.ClickRepository
}

func NewRedirectHandler(
	service *link.Service,
	clickRepositories ...*link.ClickRepository,
) *RedirectHandler {
	var clicks *link.ClickRepository

	if len(clickRepositories) > 0 {
		clicks = clickRepositories[0]
	}

	return &RedirectHandler{
		service: service,
		clicks:  clicks,
	}
}

func (h *RedirectHandler) Redirect(
	w http.ResponseWriter,
	r *http.Request,
) {
	shortCode := strings.TrimPrefix(
		r.URL.Path,
		"/",
	)

	if shortCode == "" ||
		strings.Contains(shortCode, "/") {
		http.NotFound(w, r)
		return
	}

	result, err := h.service.GetByShortCode(
		r.Context(),
		shortCode,
	)

	if errors.Is(err, link.ErrLinkNotFound) {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	now := time.Now().UTC()

	if !result.IsAvailable(now) {
		http.NotFound(w, r)
		return
	}

	if h.clicks != nil {
		event := link.ClickEvent{
			ShortCode: shortCode,
			ClickedAt: now,
			Referrer:  r.Referer(),
			UserAgent: r.UserAgent(),
			ClientIP:  clientIP(r),
		}

		// Analytics failures must never prevent
		// a valid redirect.
		_ = h.clicks.Record(
			r.Context(),
			event,
		)
	}

	http.Redirect(
		w,
		r,
		result.Destination,
		http.StatusFound,
	)
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
