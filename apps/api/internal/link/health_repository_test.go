package link

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestSQLiteHealthCheckRepositoryCreateAndGetLatest(
	t *testing.T,
) {
	db := newHealthRepositoryTestDB(t)
	repository := NewHealthCheckRepository(db)

	checkedAt := time.Date(
		2026,
		10,
		1,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	check := LinkHealth{
		ShortCode:      "abc123",
		Destination:    "https://example.com",
		Status:         HealthStatusHealthy,
		HTTPStatus:     200,
		ResponseTimeMS: 125,
		FinalURL:       "https://example.com/",
		Redirects:      1,
		HTTPS:          true,
		CheckedAt:      checkedAt,
	}

	if err := repository.Create(
		context.Background(),
		check,
	); err != nil {
		t.Fatalf(
			"create health check: %v",
			err,
		)
	}

	result, err := repository.GetLatestByShortCode(
		context.Background(),
		"abc123",
	)
	if err != nil {
		t.Fatalf(
			"get latest health check: %v",
			err,
		)
	}

	if result.ShortCode != check.ShortCode {
		t.Fatalf(
			"expected short code %q, got %q",
			check.ShortCode,
			result.ShortCode,
		)
	}

	if result.Status != check.Status {
		t.Fatalf(
			"expected status %q, got %q",
			check.Status,
			result.Status,
		)
	}

	if result.HTTPStatus != check.HTTPStatus {
		t.Fatalf(
			"expected HTTP status %d, got %d",
			check.HTTPStatus,
			result.HTTPStatus,
		)
	}

	if result.ResponseTimeMS != check.ResponseTimeMS {
		t.Fatalf(
			"expected response time %d, got %d",
			check.ResponseTimeMS,
			result.ResponseTimeMS,
		)
	}

	if result.FinalURL != check.FinalURL {
		t.Fatalf(
			"expected final URL %q, got %q",
			check.FinalURL,
			result.FinalURL,
		)
	}

	if result.Redirects != check.Redirects {
		t.Fatalf(
			"expected redirects %d, got %d",
			check.Redirects,
			result.Redirects,
		)
	}

	if result.HTTPS != check.HTTPS {
		t.Fatalf(
			"expected HTTPS %t, got %t",
			check.HTTPS,
			result.HTTPS,
		)
	}

	if !result.CheckedAt.Equal(check.CheckedAt) {
		t.Fatalf(
			"expected checked_at %s, got %s",
			check.CheckedAt,
			result.CheckedAt,
		)
	}
}

func TestSQLiteHealthCheckRepositoryList(
	t *testing.T,
) {
	db := newHealthRepositoryTestDB(t)
	repository := NewHealthCheckRepository(db)

	baseTime := time.Date(
		2026,
		10,
		1,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	checks := []LinkHealth{
		{
			ShortCode:      "abc123",
			Destination:    "https://example.com",
			Status:         HealthStatusHealthy,
			HTTPStatus:     200,
			ResponseTimeMS: 100,
			Redirects:      0,
			HTTPS:          true,
			CheckedAt:      baseTime,
		},
		{
			ShortCode:      "abc123",
			Destination:    "https://example.com",
			Status:         HealthStatusDegraded,
			HTTPStatus:     302,
			ResponseTimeMS: 500,
			FinalURL:       "https://example.com/next",
			Redirects:      1,
			HTTPS:          true,
			CheckedAt:      baseTime.Add(time.Minute),
		},
		{
			ShortCode:      "abc123",
			Destination:    "https://example.com",
			Status:         HealthStatusHealthy,
			HTTPStatus:     200,
			ResponseTimeMS: 150,
			Redirects:      0,
			HTTPS:          true,
			CheckedAt:      baseTime.Add(2 * time.Minute),
		},
		{
			ShortCode:      "other",
			Destination:    "https://example.org",
			Status:         HealthStatusHealthy,
			HTTPStatus:     200,
			ResponseTimeMS: 90,
			Redirects:      0,
			HTTPS:          true,
			CheckedAt:      baseTime.Add(3 * time.Minute),
		},
	}

	for _, check := range checks {
		if err := repository.Create(
			context.Background(),
			check,
		); err != nil {
			t.Fatalf(
				"create health check: %v",
				err,
			)
		}
	}

	results, err := repository.ListByShortCode(
		context.Background(),
		"abc123",
		2,
	)
	if err != nil {
		t.Fatalf(
			"list health checks: %v",
			err,
		)
	}

	if len(results) != 2 {
		t.Fatalf(
			"expected 2 health checks, got %d",
			len(results),
		)
	}

	if results[0].CheckedAt != baseTime.Add(2*time.Minute) {
		t.Fatalf(
			"expected newest check first, got %s",
			results[0].CheckedAt,
		)
	}

	if results[1].CheckedAt != baseTime.Add(time.Minute) {
		t.Fatalf(
			"expected second newest check, got %s",
			results[1].CheckedAt,
		)
	}
}

func TestSQLiteHealthCheckRepositoryLatestNotFound(
	t *testing.T,
) {
	db := newHealthRepositoryTestDB(t)
	repository := NewHealthCheckRepository(db)

	_, err := repository.GetLatestByShortCode(
		context.Background(),
		"does-not-exist",
	)

	if err != ErrHealthCheckNotFound {
		t.Fatalf(
			"expected ErrHealthCheckNotFound, got %v",
			err,
		)
	}
}

func newHealthRepositoryTestDB(
	t *testing.T,
) *sql.DB {
	t.Helper()

	db, err := sql.Open(
		"sqlite",
		":memory:",
	)
	if err != nil {
		t.Fatalf(
			"open test database: %v",
			err,
		)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf(
			"initialize test database: %v",
			err,
		)
	}

	return db
}
