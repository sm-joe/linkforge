package httpserver

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sm-joe/linkforge/internal/link"
)

type AnalyticsHandler struct {
	service *link.Service
	clicks  *link.ClickRepository
}

func NewAnalyticsHandler(
	service *link.Service,
	clicks *link.ClickRepository,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
		clicks:  clicks,
	}
}

type analyticsResponse struct {
	ShortCode   string              `json:"short_code"`
	TotalClicks int64               `json:"total_clicks"`
	FirstClick  *time.Time          `json:"first_click,omitempty"`
	LastClick   *time.Time          `json:"last_click,omitempty"`
	Referrers   []referrerResponse  `json:"referrers"`
	Browsers    []userAgentResponse `json:"browsers"`
	Devices     []userAgentResponse `json:"devices"`
	ClientIPs   []link.ClientIPStat `json:"client_ips"`
}

type referrerResponse struct {
	Referrer string `json:"referrer"`
	Clicks   int64  `json:"clicks"`
}

type userAgentResponse struct {
	Name   string `json:"name"`
	Clicks int64  `json:"clicks"`
}

func (h *AnalyticsHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {
	const prefix = "/api/v1/links/"

	path := strings.TrimPrefix(
		r.URL.Path,
		prefix,
	)

	if !strings.HasSuffix(path, "/analytics") {
		writeError(
			w,
			http.StatusNotFound,
			"analytics endpoint not found",
		)
		return
	}

	code := strings.TrimSuffix(
		path,
		"/analytics",
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

	summary, err := h.clicks.GetSummaryByShortCode(
		r.Context(),
		code,
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to read click summary",
		)
		return
	}

	referrers, err := h.clicks.ListReferrersByShortCode(
		r.Context(),
		code,
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to read referrers",
		)
		return
	}

	browsers, err := h.clicks.ListBrowsersByShortCode(
		r.Context(),
		code,
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to read browsers",
		)
		return
	}

	devices, err := h.clicks.ListDevicesByShortCode(
		r.Context(),
		code,
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to read devices",
		)
		return
	}

	clientIPs, err := h.clicks.ListClientIPsByShortCode(
		r.Context(),
		code,
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to read client IPs",
		)
		return
	}

	referrerResults := make(
		[]referrerResponse,
		0,
		len(referrers),
	)

	for _, stat := range referrers {
		referrerResults = append(
			referrerResults,
			referrerResponse{
				Referrer: stat.Referrer,
				Clicks:   stat.Clicks,
			},
		)
	}

	browserResults := make(
		[]userAgentResponse,
		0,
		len(browsers),
	)

	for _, stat := range browsers {
		browserResults = append(
			browserResults,
			userAgentResponse{
				Name:   stat.Name,
				Clicks: stat.Clicks,
			},
		)
	}

	deviceResults := make(
		[]userAgentResponse,
		0,
		len(devices),
	)

	for _, stat := range devices {
		deviceResults = append(
			deviceResults,
			userAgentResponse{
				Name:   stat.Name,
				Clicks: stat.Clicks,
			},
		)
	}

	writeJSON(
		w,
		http.StatusOK,
		analyticsResponse{
			ShortCode:   code,
			TotalClicks: summary.TotalClicks,
			FirstClick:  summary.FirstClick,
			LastClick:   summary.LastClick,
			Referrers:   referrerResults,
			Browsers:    browserResults,
			Devices:     deviceResults,
			ClientIPs:   clientIPs,
		},
	)
}
