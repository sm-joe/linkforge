package link

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const maxHealthRedirects = 5

type HealthService struct {
	client     *HealthClient
	repository HealthCheckRepository
}

func NewHealthService(
	client *HealthClient,
	repositories ...HealthCheckRepository,
) *HealthService {
	var repository HealthCheckRepository

	if len(repositories) > 0 {
		repository = repositories[0]
	}

	return &HealthService{
		client:     client,
		repository: repository,
	}
}

func (s *HealthService) Check(
	ctx context.Context,
	link *Link,
) LinkHealth {
	start := time.Now()

	result := LinkHealth{
		ShortCode:   link.ShortCode,
		Destination: link.Destination,
		HTTPS:       false,
		CheckedAt:   start.UTC(),
	}

	currentURL := link.Destination
	redirects := 0

	for {
		response, err := s.client.Get(
			ctx,
			currentURL,
		)

		if err != nil {
			result.Status = HealthStatusUnhealthy
			result.Error = err.Error()
			result.Redirects = redirects
			result.ResponseTimeMS =
				time.Since(start).Milliseconds()

			s.persist(ctx, result)

			return result
		}

		result.HTTPStatus = response.StatusCode
		result.FinalURL = currentURL
		result.HTTPS = isHTTPSURL(currentURL)

		if !isRedirectStatus(response.StatusCode) {
			_ = response.Body.Close()
			break
		}

		location := response.Header.Get("Location")

		_ = response.Body.Close()

		if location == "" {
			result.Status = HealthStatusUnhealthy
			result.Error =
				"redirect response missing Location header"
			result.Redirects = redirects
			result.ResponseTimeMS =
				time.Since(start).Milliseconds()

			s.persist(ctx, result)

			return result
		}

		redirects++

		if redirects > maxHealthRedirects {
			result.Status = HealthStatusUnhealthy
			result.Error = fmt.Sprintf(
				"redirect limit exceeded (%d)",
				maxHealthRedirects,
			)
			result.Redirects = redirects
			result.ResponseTimeMS =
				time.Since(start).Milliseconds()

			s.persist(ctx, result)

			return result
		}

		nextURL, err := resolveRedirectURL(
			currentURL,
			location,
		)
		if err != nil {
			result.Status = HealthStatusUnhealthy
			result.Error = err.Error()
			result.Redirects = redirects
			result.ResponseTimeMS =
				time.Since(start).Milliseconds()

			s.persist(ctx, result)

			return result
		}

		currentURL = nextURL
	}

	result.Redirects = redirects
	result.ResponseTimeMS =
		time.Since(start).Milliseconds()

	switch {
	case result.HTTPStatus >= http.StatusOK &&
		result.HTTPStatus < http.StatusMultipleChoices:

		result.Status = HealthStatusHealthy

	case result.HTTPStatus >= http.StatusMultipleChoices &&
		result.HTTPStatus < http.StatusBadRequest:

		result.Status = HealthStatusDegraded
		result.Error = fmt.Sprintf(
			"destination returned HTTP %d",
			result.HTTPStatus,
		)

	default:
		result.Status = HealthStatusUnhealthy
		result.Error = fmt.Sprintf(
			"destination returned HTTP %d",
			result.HTTPStatus,
		)
	}

	s.persist(ctx, result)

	return result
}

func (s *HealthService) persist(
	ctx context.Context,
	result LinkHealth,
) {
	if s.repository == nil {
		return
	}

	_ = s.repository.Create(
		ctx,
		result,
	)
}

func isRedirectStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect:
		return true

	default:
		return false
	}
}

func resolveRedirectURL(
	currentURL string,
	location string,
) (string, error) {
	base, err := url.Parse(currentURL)
	if err != nil {
		return "", fmt.Errorf(
			"parse current URL: %w",
			err,
		)
	}

	target, err := url.Parse(location)
	if err != nil {
		return "", fmt.Errorf(
			"parse redirect location: %w",
			err,
		)
	}

	resolved := base.ResolveReference(target)

	if resolved.Scheme != "http" &&
		resolved.Scheme != "https" {
		return "", fmt.Errorf(
			"redirect uses unsupported scheme",
		)
	}

	if resolved.Hostname() == "" {
		return "", fmt.Errorf(
			"redirect has no hostname",
		)
	}

	return resolved.String(), nil
}

func isHTTPSURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return parsed.Scheme == "https"
}
