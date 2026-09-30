package httpserver

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sm-joe/linkforge/internal/link"
)

type RedirectHandler struct {
	service *link.Service
}

func NewRedirectHandler(service *link.Service) *RedirectHandler {
	return &RedirectHandler{
		service: service,
	}
}

func (h *RedirectHandler) Redirect(
	w http.ResponseWriter,
	r *http.Request,
) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
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

	http.Redirect(
		w,
		r,
		result.Destination,
		http.StatusFound,
	)
}
