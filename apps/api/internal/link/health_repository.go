package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrHealthCheckNotFound = errors.New(
	"health check not found",
)

type HealthCheckRepository interface {
	Create(
		ctx context.Context,
		check LinkHealth,
	) error

	ListByShortCode(
		ctx context.Context,
		shortCode string,
		limit int,
	) ([]LinkHealth, error)

	GetLatestByShortCode(
		ctx context.Context,
		shortCode string,
	) (*LinkHealth, error)
}

type SQLiteHealthCheckRepository struct {
	db *sql.DB
}

func NewHealthCheckRepository(
	db *sql.DB,
) *SQLiteHealthCheckRepository {
	return &SQLiteHealthCheckRepository{
		db: db,
	}
}

func (r *SQLiteHealthCheckRepository) Create(
	ctx context.Context,
	check LinkHealth,
) error {
	const query = `
INSERT INTO link_health_checks (
    short_code,
    destination,
    status,
    http_status,
    response_time_ms,
    final_url,
    redirects,
    https,
    checked_at,
    error
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	_, err := r.db.ExecContext(
		ctx,
		query,
		check.ShortCode,
		check.Destination,
		check.Status,
		nullableHTTPStatus(check.HTTPStatus),
		check.ResponseTimeMS,
		nullableString(check.FinalURL),
		check.Redirects,
		check.HTTPS,
		check.CheckedAt.UTC(),
		nullableString(check.Error),
	)
	if err != nil {
		return fmt.Errorf(
			"create health check: %w",
			err,
		)
	}

	return nil
}

func (r *SQLiteHealthCheckRepository) ListByShortCode(
	ctx context.Context,
	shortCode string,
	limit int,
) ([]LinkHealth, error) {
	if limit <= 0 {
		limit = 50
	}

	const query = `
SELECT
    short_code,
    destination,
    status,
    http_status,
    response_time_ms,
    final_url,
    redirects,
    https,
    checked_at,
    error
FROM link_health_checks
WHERE short_code = ?
ORDER BY checked_at DESC
LIMIT ?
`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		shortCode,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list health checks: %w",
			err,
		)
	}
	defer rows.Close()

	checks := make(
		[]LinkHealth,
		0,
	)

	for rows.Next() {
		check, err := scanHealthCheck(rows)
		if err != nil {
			return nil, err
		}

		checks = append(
			checks,
			check,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate health checks: %w",
			err,
		)
	}

	return checks, nil
}

func (r *SQLiteHealthCheckRepository) GetLatestByShortCode(
	ctx context.Context,
	shortCode string,
) (*LinkHealth, error) {
	const query = `
SELECT
    short_code,
    destination,
    status,
    http_status,
    response_time_ms,
    final_url,
    redirects,
    https,
    checked_at,
    error
FROM link_health_checks
WHERE short_code = ?
ORDER BY checked_at DESC
LIMIT 1
`

	row := r.db.QueryRowContext(
		ctx,
		query,
		shortCode,
	)

	check, err := scanHealthCheckRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrHealthCheckNotFound
		}

		return nil, fmt.Errorf(
			"get latest health check: %w",
			err,
		)
	}

	return &check, nil
}

type healthCheckScanner interface {
	Scan(dest ...any) error
}

func scanHealthCheck(
	scanner healthCheckScanner,
) (LinkHealth, error) {
	var (
		check        LinkHealth
		status       string
		httpStatus   sql.NullInt64
		finalURL     sql.NullString
		https        bool
		checkedAt    time.Time
		errorMessage sql.NullString
	)

	err := scanner.Scan(
		&check.ShortCode,
		&check.Destination,
		&status,
		&httpStatus,
		&check.ResponseTimeMS,
		&finalURL,
		&check.Redirects,
		&https,
		&checkedAt,
		&errorMessage,
	)
	if err != nil {
		return LinkHealth{}, fmt.Errorf(
			"scan health check: %w",
			err,
		)
	}

	check.Status = HealthStatus(status)
	check.HTTPStatus = int(httpStatus.Int64)
	check.FinalURL = finalURL.String
	check.HTTPS = https
	check.CheckedAt = checkedAt.UTC()
	check.Error = errorMessage.String

	if !httpStatus.Valid {
		check.HTTPStatus = 0
	}

	if !finalURL.Valid {
		check.FinalURL = ""
	}

	if !errorMessage.Valid {
		check.Error = ""
	}

	return check, nil
}

func scanHealthCheckRow(
	row *sql.Row,
) (LinkHealth, error) {
	return scanHealthCheck(row)
}

func nullableHTTPStatus(
	status int,
) any {
	if status == 0 {
		return nil
	}

	return status
}

func nullableString(
	value string,
) any {
	if value == "" {
		return nil
	}

	return value
}
