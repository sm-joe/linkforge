package httpserver

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sm-joe/linkforge/internal/link"
)

type LinkManagementHandler struct {
	service *link.Service
}

func NewLinkManagementHandler(
	service *link.Service,
) *LinkManagementHandler {
	return &LinkManagementHandler{
		service: service,
	}
}

func (h *LinkManagementHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {
	code := extractCode(
		r.URL.Path,
		"/api/v1/links/",
	)

	if code == "" {
		writeError(
			w,
			http.StatusNotFound,
			"link not found",
		)
		return
	}

	result, err := h.service.GetByShortCode(
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

	writeJSON(
		w,
		http.StatusOK,
		result,
	)
}

func (h *LinkManagementHandler) Disable(
	w http.ResponseWriter,
	r *http.Request,
) {
	code := extractCode(
		strings.TrimSuffix(
			r.URL.Path,
			"/disable",
		),
		"/api/v1/links/",
	)

	if code == "" {
		writeError(
			w,
			http.StatusNotFound,
			"link not found",
		)
		return
	}

	err := h.service.DisableLink(
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

	w.WriteHeader(http.StatusNoContent)
}

func (h *LinkManagementHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	code := extractCode(
		r.URL.Path,
		"/api/v1/links/",
	)

	if code == "" {
		writeError(
			w,
			http.StatusNotFound,
			"link not found",
		)
		return
	}

	err := h.service.DeleteLink(
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

	w.WriteHeader(http.StatusNoContent)
}

func extractCode(
	path string,
	prefix string,
) string {
	code := strings.TrimPrefix(
		path,
		prefix,
	)

	if code == "" || strings.Contains(code, "/") {
		return ""
	}

	return code
}

func (h *LinkManagementHandler) Enable(
	w http.ResponseWriter,
	r *http.Request,
) {
	code := extractCode(
		strings.TrimSuffix(
			r.URL.Path,
			"/enable",
		),
		"/api/v1/links/",
	)

	if code == "" {
		writeError(
			w,
			http.StatusNotFound,
			"link not found",
		)
		return
	}

	err := h.service.EnableLink(
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

	w.WriteHeader(http.StatusNoContent)
}
