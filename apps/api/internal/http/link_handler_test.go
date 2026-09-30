package httpserver

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sm-joe/linkforge/internal/link"
	_ "modernc.org/sqlite"
)

func newTestLinkServer(t *testing.T) *httptest.Server {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	if _, err := db.Exec(`
		CREATE TABLE links (
			id TEXT PRIMARY KEY,
			short_code TEXT NOT NULL UNIQUE,
			alias TEXT,
			destination TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			expires_at DATETIME
		);

		CREATE INDEX idx_links_short_code
			ON links(short_code);

		CREATE INDEX idx_links_status
			ON links(status);
	`); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	repository := link.NewSQLiteRepository(db)
	service := link.NewService(repository)
	handler := NewLinkHandler(service)
	redirectHandler := NewRedirectHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/links", handler.Create)
	mux.HandleFunc("/", redirectHandler.Redirect)

	return httptest.NewServer(mux)
}

func TestCreateLinkAPI(t *testing.T) {
	server := newTestLinkServer(t)
	defer server.Close()

	payload := map[string]string{
		"destination": "https://example.com",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	response, err := http.Post(
		server.URL+"/api/v1/links",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("POST request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf(
			"status = %d, want %d",
			response.StatusCode,
			http.StatusCreated,
		)
	}

	var result createLinkResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.ID == "" {
		t.Fatal("expected link ID")
	}

	if result.ShortCode == "" {
		t.Fatal("expected short code")
	}

	if result.Destination != "https://example.com" {
		t.Fatalf(
			"destination = %q, want %q",
			result.Destination,
			"https://example.com",
		)
	}

	if result.Status != link.StatusActive {
		t.Fatalf(
			"status = %q, want %q",
			result.Status,
			link.StatusActive,
		)
	}
}

func TestCreateLinkRejectsUnknownFields(t *testing.T) {
	server := newTestLinkServer(t)
	defer server.Close()

	payload := `{
		"destination": "https://example.com",
		"unknown": "field"
	}`

	response, err := http.Post(
		server.URL+"/api/v1/links",
		"application/json",
		bytes.NewBufferString(payload),
	)
	if err != nil {
		t.Fatalf("POST request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want %d",
			response.StatusCode,
			http.StatusBadRequest,
		)
	}
}

func TestCreateLinkRejectsWrongContentType(t *testing.T) {
	server := newTestLinkServer(t)
	defer server.Close()

	request, err := http.NewRequest(
		http.MethodPost,
		server.URL+"/api/v1/links",
		bytes.NewBufferString(`{"destination":"https://example.com"}`),
	)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	request.Header.Set("Content-Type", "text/plain")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("POST request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf(
			"status = %d, want %d",
			response.StatusCode,
			http.StatusUnsupportedMediaType,
		)
	}
}

func TestRedirectLink(t *testing.T) {
	server := newTestLinkServer(t)
	defer server.Close()

	payload := `{
		"destination": "https://example.com",
		"alias": "example"
	}`

	createResponse, err := http.Post(
		server.URL+"/api/v1/links",
		"application/json",
		bytes.NewBufferString(payload),
	)
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	defer createResponse.Body.Close()

	if createResponse.StatusCode != http.StatusCreated {
		t.Fatalf(
			"create status = %d, want %d",
			createResponse.StatusCode,
			http.StatusCreated,
		)
	}

	client := &http.Client{
		CheckRedirect: func(
			req *http.Request,
			via []*http.Request,
		) error {
			return http.ErrUseLastResponse
		},
	}

	redirectResponse, err := client.Get(
		server.URL + "/example",
	)
	if err != nil {
		t.Fatalf("redirect request: %v", err)
	}
	defer redirectResponse.Body.Close()

	if redirectResponse.StatusCode != http.StatusFound {
		t.Fatalf(
			"redirect status = %d, want %d",
			redirectResponse.StatusCode,
			http.StatusFound,
		)
	}

	location := redirectResponse.Header.Get("Location")

	if location != "https://example.com" {
		t.Fatalf(
			"Location = %q, want %q",
			location,
			"https://example.com",
		)
	}
}

func TestRedirectMissingLink(t *testing.T) {
	server := newTestLinkServer(t)
	defer server.Close()

	response, err := http.Get(
		server.URL + "/does-not-exist",
	)
	if err != nil {
		t.Fatalf("GET request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"status = %d, want %d",
			response.StatusCode,
			http.StatusNotFound,
		)
	}
}
