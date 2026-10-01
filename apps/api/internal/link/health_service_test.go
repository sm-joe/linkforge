package link

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	urlvalidation "github.com/sm-joe/linkforge/packages/url-validation"
)

type healthTestResolver struct {
	ips []net.IP
	err error
}

func (r healthTestResolver) LookupIP(
	ctx context.Context,
	network string,
	host string,
) ([]net.IP, error) {
	return r.ips, r.err
}

func TestHealthServiceCheckHealthy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestHTTPClient(server)
	service := NewHealthService(client)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "test123",
			Destination: server.URL,
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusHealthy {
		t.Fatalf(
			"expected healthy, got %s",
			result.Status,
		)
	}

	if result.HTTPStatus != http.StatusOK {
		t.Fatalf(
			"expected HTTP 200, got %d",
			result.HTTPStatus,
		)
	}

	if result.ResponseTimeMS < 0 {
		t.Fatalf(
			"expected non-negative response time, got %d",
			result.ResponseTimeMS,
		)
	}
}

func TestHealthServiceCheckUnhealthyHTTPStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestHTTPClient(server)
	service := NewHealthService(client)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "test123",
			Destination: server.URL,
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusUnhealthy {
		t.Fatalf(
			"expected unhealthy, got %s",
			result.Status,
		)
	}

	if result.HTTPStatus != http.StatusNotFound {
		t.Fatalf(
			"expected HTTP 404, got %d",
			result.HTTPStatus,
		)
	}

	if result.Error == "" {
		t.Fatal("expected health error")
	}
}

func TestHealthServiceCheckPrivateDestination(t *testing.T) {
	resolver := healthTestResolver{
		ips: []net.IP{
			net.ParseIP("10.0.0.10"),
		},
	}

	client := NewHealthClient(resolver)
	service := NewHealthService(client)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "private",
			Destination: "http://example.com",
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusUnhealthy {
		t.Fatalf(
			"expected unhealthy, got %s",
			result.Status,
		)
	}

	if result.Error == "" {
		t.Fatal("expected private destination error")
	}
}

func TestHealthServiceCheckDNSFailure(t *testing.T) {
	resolver := healthTestResolver{
		err: errors.New("DNS unavailable"),
	}

	client := NewHealthClient(resolver)
	service := NewHealthService(client)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "dns-failure",
			Destination: "https://example.com",
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusUnhealthy {
		t.Fatalf(
			"expected unhealthy, got %s",
			result.Status,
		)
	}

	if result.Error == "" {
		t.Fatal("expected DNS error")
	}
}

func TestHealthServiceCheckRedirectChain(t *testing.T) {
	finalServer := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
	}))
	defer finalServer.Close()

	redirectServer := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		http.Redirect(
			w,
			r,
			finalServer.URL,
			http.StatusFound,
		)
	}))
	defer redirectServer.Close()

	client := newTestHTTPClient(redirectServer)
	service := NewHealthService(client)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "redirect",
			Destination: redirectServer.URL,
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusHealthy {
		t.Fatalf(
			"expected healthy, got %s",
			result.Status,
		)
	}

	if result.HTTPStatus != http.StatusOK {
		t.Fatalf(
			"expected HTTP 200, got %d",
			result.HTTPStatus,
		)
	}

	if result.Redirects != 1 {
		t.Fatalf(
			"expected 1 redirect, got %d",
			result.Redirects,
		)
	}

	if result.FinalURL != finalServer.URL {
		t.Fatalf(
			"expected final URL %q, got %q",
			finalServer.URL,
			result.FinalURL,
		)
	}
}

func TestHealthServiceCheckPrivateRedirectTarget(t *testing.T) {
	redirectServer := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Header().Set(
			"Location",
			"http://example.com/private",
		)
		w.WriteHeader(http.StatusFound)
	}))
	defer redirectServer.Close()

	resolver := healthTestResolver{
		ips: []net.IP{
			net.ParseIP("10.0.0.10"),
		},
	}

	productionClient := NewHealthClient(resolver)

	// Use the real validation path, but provide a transport that
	// returns the local redirect response without making an outbound
	// request to the test server through DNS.
	productionClient.client.Transport = roundTripFunc(
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusFound,
				Header: http.Header{
					"Location": []string{
						"http://example.com/private",
					},
				},
				Body: http.NoBody,
				Request: req,
			}, nil
		},
	)

	service := NewHealthService(productionClient)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "private-redirect",
			Destination: "http://example.com",
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusUnhealthy {
		t.Fatalf(
			"expected unhealthy, got %s",
			result.Status,
		)
	}

	if result.Error == "" {
		t.Fatal("expected private redirect error")
	}
}

func TestHealthServiceCheckDNSFailureOnRedirectTarget(t *testing.T) {
	resolver := healthTestResolver{
		ips: []net.IP{
			net.ParseIP("8.8.8.8"),
		},
	}

	client := NewHealthClient(resolver)

	client.client.Transport = roundTripFunc(
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusFound,
				Header: http.Header{
					"Location": []string{
						"https://redirect-target.example",
					},
				},
				Body: http.NoBody,
				Request: req,
			}, nil
		},
	)

	// The resolver is changed after the first request so the redirect
	// target fails DNS validation.
	dynamicResolver := &sequenceResolver{
		results: []resolverResult{
			{
				ips: []net.IP{
					net.ParseIP("8.8.8.8"),
				},
			},
			{
				err: errors.New("redirect DNS unavailable"),
			},
		},
	}

	client = NewHealthClient(dynamicResolver)
	client.client.Transport = roundTripFunc(
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusFound,
				Header: http.Header{
					"Location": []string{
						"https://redirect-target.example",
					},
				},
				Body: http.NoBody,
				Request: req,
			}, nil
		},
	)

	service := NewHealthService(client)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "redirect-dns",
			Destination: "https://example.com",
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusUnhealthy {
		t.Fatalf(
			"expected unhealthy, got %s",
			result.Status,
		)
	}

	if result.Error == "" {
		t.Fatal("expected redirect DNS error")
	}
}

func TestHealthServiceCheckRedirectLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		next := "http://" + r.Host + r.URL.Path

		http.Redirect(
			w,
			r,
			next,
			http.StatusFound,
		)
	}))
	defer server.Close()

	client := newTestHTTPClient(server)
	service := NewHealthService(client)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "redirect-loop",
			Destination: server.URL,
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusUnhealthy {
		t.Fatalf(
			"expected unhealthy, got %s",
			result.Status,
		)
	}

	if result.Redirects <= maxHealthRedirects {
		t.Fatalf(
			"expected redirects > %d, got %d",
			maxHealthRedirects,
			result.Redirects,
		)
	}

	if result.Error == "" {
		t.Fatal("expected redirect limit error")
	}
}

func TestHealthServiceCheckMissingLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusFound)
	}))
	defer server.Close()

	client := newTestHTTPClient(server)
	service := NewHealthService(client)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "missing-location",
			Destination: server.URL,
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusUnhealthy {
		t.Fatalf(
			"expected unhealthy, got %s",
			result.Status,
		)
	}

	if result.Error != "redirect response missing Location header" {
		t.Fatalf(
			"unexpected error: %s",
			result.Error,
		)
	}
}

func TestHealthServiceCheckRelativeRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.URL.Path == "/start" {
			http.Redirect(
				w,
				r,
				"/final",
				http.StatusFound,
			)
			return
		}

		if r.URL.Path == "/final" {
			w.WriteHeader(http.StatusOK)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client := newTestHTTPClient(server)
	service := NewHealthService(client)

	result := service.Check(
		context.Background(),
		&Link{
			ShortCode:   "relative",
			Destination: server.URL + "/start",
			Status:      StatusActive,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if result.Status != HealthStatusHealthy {
		t.Fatalf(
			"expected healthy, got %s",
			result.Status,
		)
	}

	if result.Redirects != 1 {
		t.Fatalf(
			"expected 1 redirect, got %d",
			result.Redirects,
		)
	}

	expectedFinalURL := server.URL + "/final"

	if result.FinalURL != expectedFinalURL {
		t.Fatalf(
			"expected final URL %q, got %q",
			expectedFinalURL,
			result.FinalURL,
		)
	}
}

func newTestHTTPClient(server *httptest.Server) *HealthClient {
	httpClient := server.Client()

	httpClient.CheckRedirect = func(
		req *http.Request,
		via []*http.Request,
	) error {
		return http.ErrUseLastResponse
	}

	return NewHealthClientWithHTTPClient(httpClient)
}

type roundTripFunc func(
	req *http.Request,
) (*http.Response, error)

func (f roundTripFunc) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	return f(req)
}

type resolverResult struct {
	ips []net.IP
	err error
}

type sequenceResolver struct {
	results []resolverResult
	index   int
}

func (r *sequenceResolver) LookupIP(
	ctx context.Context,
	network string,
	host string,
) ([]net.IP, error) {
	if r.index >= len(r.results) {
		return nil, errors.New("unexpected DNS lookup")
	}

	result := r.results[r.index]
	r.index++

	return result.ips, result.err
}

var _ urlvalidation.Resolver = healthTestResolver{}