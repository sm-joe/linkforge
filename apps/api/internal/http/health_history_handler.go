package httpserver

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/sm-joe/linkforge/internal/link"
)

const defaultHealthHistoryLimit = 50

type HealthHistoryHandler struct {
	service    *link.Service
	repository link.HealthCheckRepository
}

func NewHealthHistoryHandler(
	service *link.Service,
	repository link.HealthCheckRepository,
) *HealthHistoryHandler {
	return &HealthHistoryHandler{
		service:    service,
		repository: repository,
	}
}

type healthHistoryResponse struct {
	ShortCode string            `json:"short_code"`
	Checks    []link.LinkHealth `json:"checks"`
}

func (h *HealthHistoryHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {
	const prefix = "/api/v1/links/"

	path := strings.TrimPrefix(
		r.URL.Path,
		prefix,
	)

	if !strings.HasSuffix(path, "/health/history") {
		writeError(
			w,
			http.StatusNotFound,
			"health history endpoint not found",
		)
		return
	}

	code := strings.TrimSuffix(
		path,
		"/health/history",
	)

	code = strings.TrimSpace(code)

	if code == "" || strings.Contains(code, "/") {
		writeError(
			w,
			http.StatusNotFound,
			"link not found",
		)
		return
	}

	_, err := h.service.GetByShortCode(
		r.Context(),
		code,
	)

	if errors.Is(err, link.ErrLinkNotFound) {
		writeError(
			w,
			http.StatusNotFound,
			"link not found",
		)
		return
	}

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	limit := defaultHealthHistoryLimit

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			writeError(
				w,
				http.StatusBadRequest,
				"invalid limit",
			)
			return
		}

		if parsed < limit {
			limit = parsed
		}
	}

	checks, err := h.repository.ListByShortCode(
		r.Context(),
		code,
		limit,
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to read health history",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		healthHistoryResponse{
			ShortCode: code,
			Checks:    checks,
		},
	)
}
