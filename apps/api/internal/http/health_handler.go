package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sm-joe/linkforge/internal/link"
)

type HealthHandler struct {
	linkService   *link.Service
	healthService *link.HealthService
}

func NewHealthHandler(
	linkService *link.Service,
	healthService *link.HealthService,
) *HealthHandler {
	return &HealthHandler{
		linkService:   linkService,
		healthService: healthService,
	}
}

func (h *HealthHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {
	shortCode := extractHealthShortCode(r.URL.Path)

	if shortCode == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"short code is required",
		)
		return
	}

	managedLink, err := h.linkService.GetByShortCode(
		r.Context(),
		shortCode,
	)
	if err != nil {
		if err == link.ErrLinkNotFound {
			writeError(
				w,
				http.StatusNotFound,
				"link not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to retrieve link",
		)
		return
	}

	result := h.healthService.Check(
		r.Context(),
		managedLink,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(result)
}

func extractHealthShortCode(path string) string {
	const prefix = "/api/v1/links/"
	const suffix = "/health"

	if !strings.HasPrefix(path, prefix) ||
		!strings.HasSuffix(path, suffix) {
		return ""
	}

	shortCode := strings.TrimPrefix(path, prefix)
	shortCode = strings.TrimSuffix(shortCode, suffix)
	shortCode = strings.Trim(shortCode, "/")

	if strings.Contains(shortCode, "/") {
		return ""
	}

	return shortCode
}
