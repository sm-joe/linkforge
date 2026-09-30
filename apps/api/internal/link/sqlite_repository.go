package link

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"modernc.org/sqlite"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Create(
	ctx context.Context,
	link *Link,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO links (
			id,
			short_code,
			alias,
			destination,
			status,
			created_at,
			expires_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
		link.ID,
		link.ShortCode,
		link.Alias,
		link.Destination,
		link.Status,
		link.CreatedAt.UTC(),
		link.ExpiresAt,
	)

	if err != nil {
		var sqliteErr *sqlite.Error

		if errors.As(err, &sqliteErr) {
			if strings.Contains(
				strings.ToLower(sqliteErr.Error()),
				"unique",
			) {
				return ErrShortCodeTaken
			}
		}

		return err
	}

	return nil
}

func (r *SQLiteRepository) GetByShortCode(
	ctx context.Context,
	shortCode string,
) (*Link, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			id,
			short_code,
			alias,
			destination,
			status,
			created_at,
			expires_at
		FROM links
		WHERE short_code = ?
		`,
		shortCode,
	)

	var link Link
	var expiresAt sql.NullTime

	err := row.Scan(
		&link.ID,
		&link.ShortCode,
		&link.Alias,
		&link.Destination,
		&link.Status,
		&link.CreatedAt,
		&expiresAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLinkNotFound
	}

	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		expires := expiresAt.Time
		link.ExpiresAt = &expires
	}

	return &link, nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

var _ Repository = (*SQLiteRepository)(nil)

func (r *SQLiteRepository) UpdateStatus(
	ctx context.Context,
	shortCode string,
	status Status,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
		UPDATE links
		SET status = ?
		WHERE short_code = ?
		`,
		status,
		shortCode,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrLinkNotFound
	}

	return nil
}

func (r *SQLiteRepository) Delete(
	ctx context.Context,
	shortCode string,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
		DELETE FROM links
		WHERE short_code = ?
		`,
		shortCode,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrLinkNotFound
	}

	return nil
}

func (r *SQLiteRepository) List(
	ctx context.Context,
) ([]*Link, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			id,
			short_code,
			alias,
			destination,
			status,
			created_at,
			expires_at
		FROM links
		ORDER BY created_at DESC
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var links []*Link

	for rows.Next() {
		var link Link
		var expiresAt sql.NullTime

		err := rows.Scan(
			&link.ID,
			&link.ShortCode,
			&link.Alias,
			&link.Destination,
			&link.Status,
			&link.CreatedAt,
			&expiresAt,
		)

		if err != nil {
			return nil, err
		}

		if expiresAt.Valid {
			value := expiresAt.Time
			link.ExpiresAt = &value
		}

		links = append(
			links,
			&link,
		)
	}

	return links, rows.Err()
}
