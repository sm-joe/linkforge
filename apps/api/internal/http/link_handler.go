package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sm-joe/linkforge/internal/link"
)

type LinkHandler struct {
	service *link.Service
}

func NewLinkHandler(service *link.Service) *LinkHandler {
	return &LinkHandler{
		service: service,
	}
}

type createLinkRequest struct {
	Destination string     `json:"destination"`
	Alias       string     `json:"alias,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type createLinkResponse struct {
	ID          string      `json:"id"`
	ShortCode   string      `json:"short_code"`
	ShortURL    string      `json:"short_url"`
	Destination string      `json:"destination"`
	Status      link.Status `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
}

func (h *LinkHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	contentType := r.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "application/json") {
		writeError(
			w,
			http.StatusUnsupportedMediaType,
			"Content-Type must be application/json",
		)
		return
	}

	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		1<<20,
	)
	defer r.Body.Close()

	var request createLinkRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid JSON request",
		)
		return
	}

	result, err := h.service.CreateLink(
		r.Context(),
		link.CreateRequest{
			Destination: request.Destination,
			Alias:       request.Alias,
			ExpiresAt:   request.ExpiresAt,
		},
	)

	if err != nil {
		status := http.StatusInternalServerError

		switch {
		case errors.Is(err, link.ErrDestinationRequired),
			errors.Is(err, link.ErrInvalidAlias),
			errors.Is(err, link.ErrAliasTooLong):
			status = http.StatusBadRequest

		case errors.Is(err, link.ErrShortCodeTaken):
			status = http.StatusConflict
		}

		writeError(w, status, err.Error())
		return
	}

	response := createLinkResponse{
		ID:          result.ID,
		ShortCode:   result.ShortCode,
		ShortURL:    buildShortURL(r, result.ShortCode),
		Destination: result.Destination,
		Status:      result.Status,
		CreatedAt:   result.CreatedAt,
		ExpiresAt:   result.ExpiresAt,
	}

	writeJSON(
		w,
		http.StatusCreated,
		response,
	)
}

func buildShortURL(
	r *http.Request,
	shortCode string,
) string {
	scheme := "http"

	if r.TLS != nil {
		scheme = "https"
	}

	host := r.Host

	if strings.TrimSpace(host) == "" {
		host = "localhost:8080"
	}

	return scheme + "://" + host + "/" + shortCode
}
