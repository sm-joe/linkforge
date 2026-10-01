package httpserver

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sm-joe/linkforge/internal/link"
	"github.com/skip2/go-qrcode"
)

type QRHandler struct {
	service *link.Service
}

func NewQRHandler(
	service *link.Service,
) *QRHandler {
	return &QRHandler{
		service: service,
	}
}

func (h *QRHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {
	const prefix = "/api/v1/links/"

	path := strings.TrimPrefix(
		r.URL.Path,
		prefix,
	)

	if !strings.HasSuffix(path, "/qr") {
		writeError(
			w,
			http.StatusNotFound,
			"QR endpoint not found",
		)
		return
	}

	code := strings.TrimSuffix(
		path,
		"/qr",
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

	if !result.IsAvailable(
		time.Now().UTC(),
	) {
		writeError(
			w,
			http.StatusNotFound,
			"link not available",
		)
		return
	}

	shortURL := buildShortURL(
		r,
		result.ShortCode,
	)

	png, err := qrcode.Encode(
		shortURL,
		qrcode.Medium,
		256,
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to generate QR code",
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"image/png",
	)

	w.Header().Set(
		"Cache-Control",
		"public, max-age=300",
	)

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(png)
}