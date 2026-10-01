package link

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ClickEvent struct {
	ID        int64     `json:"id"`
	ShortCode string    `json:"short_code"`
	ClickedAt time.Time `json:"clicked_at"`
	Referrer  string    `json:"referrer,omitempty"`
	UserAgent string    `json:"user_agent"`
	ClientIP  string    `json:"client_ip,omitempty"`
}

type ClickRepository struct {
	db *sql.DB
}

func NewClickRepository(db *sql.DB) *ClickRepository {
	return &ClickRepository{
		db: db,
	}
}

type ClickSummary struct {
	TotalClicks int64
	FirstClick  *time.Time
	LastClick   *time.Time
}

type ReferrerStat struct {
	Referrer string
	Clicks   int64
}

type UserAgentStat struct {
	Name   string
	Clicks int64
}

type ClientIPStat struct {
	IP     string `json:"ip"`
	Clicks int64  `json:"clicks"`
}

func (r *ClickRepository) Record(
	ctx context.Context,
	event ClickEvent,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO click_events (
			short_code,
			clicked_at,
			referrer,
			user_agent,
			client_ip
		)
		VALUES (?, ?, ?, ?, ?)
		`,
		event.ShortCode,
		event.ClickedAt.UTC(),
		event.Referrer,
		event.UserAgent,
		event.ClientIP,
	)

	if err != nil {
		return fmt.Errorf(
			"record click event: %w",
			err,
		)
	}

	return nil
}

func (r *ClickRepository) CountByShortCode(
	ctx context.Context,
	shortCode string,
) (int64, error) {
	var count int64

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT COUNT(*)
		FROM click_events
		WHERE short_code = ?
		`,
		shortCode,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf(
			"count clicks: %w",
			err,
		)
	}

	return count, nil
}

func (r *ClickRepository) GetSummaryByShortCode(
	ctx context.Context,
	shortCode string,
) (ClickSummary, error) {
	var summary ClickSummary

	var firstClick string
	var lastClick string

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			COUNT(*),
			MIN(clicked_at),
			MAX(clicked_at)
		FROM click_events
		WHERE short_code = ?
		`,
		shortCode,
	).Scan(
		&summary.TotalClicks,
		&firstClick,
		&lastClick,
	)

	if err != nil {
		return ClickSummary{}, fmt.Errorf(
			"get click summary: %w",
			err,
		)
	}

	if firstClick != "" {
		value, err := parseSQLiteTime(firstClick)
		if err != nil {
			return ClickSummary{}, fmt.Errorf(
				"parse first click time: %w",
				err,
			)
		}

		summary.FirstClick = &value
	}

	if lastClick != "" {
		value, err := parseSQLiteTime(lastClick)
		if err != nil {
			return ClickSummary{}, fmt.Errorf(
				"parse last click time: %w",
				err,
			)
		}

		summary.LastClick = &value
	}

	return summary, nil
}

func parseSQLiteTime(value string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05.999999 -0700 MST",
		"2006-01-02 15:04:05 -0700 MST",
	}

	var lastErr error

	for _, format := range formats {
		parsed, err := time.Parse(format, value)
		if err == nil {
			return parsed.UTC(), nil
		}

		lastErr = err
	}

	return time.Time{}, lastErr
}

func (r *ClickRepository) ListReferrersByShortCode(
	ctx context.Context,
	shortCode string,
) ([]ReferrerStat, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			CASE
				WHEN TRIM(referrer) = '' THEN '(direct)'
				ELSE referrer
			END AS referrer,
			COUNT(*) AS clicks
		FROM click_events
		WHERE short_code = ?
		GROUP BY
			CASE
				WHEN TRIM(referrer) = '' THEN '(direct)'
				ELSE referrer
			END
		ORDER BY clicks DESC, referrer ASC
		`,
		shortCode,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list referrers: %w",
			err,
		)
	}
	defer rows.Close()

	referrers := make([]ReferrerStat, 0)

	for rows.Next() {
		var stat ReferrerStat

		if err := rows.Scan(
			&stat.Referrer,
			&stat.Clicks,
		); err != nil {
			return nil, fmt.Errorf(
				"scan referrer: %w",
				err,
			)
		}

		referrers = append(
			referrers,
			stat,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate referrers: %w",
			err,
		)
	}

	return referrers, nil
}

func (r *ClickRepository) ListBrowsersByShortCode(
	ctx context.Context,
	shortCode string,
) ([]UserAgentStat, error) {
	events, err := r.ListByShortCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf(
			"list browsers: %w",
			err,
		)
	}

	counts := make(map[string]int64)

	for _, event := range events {
		browser := DetectBrowser(event.UserAgent)
		counts[browser]++
	}

	stats := make([]UserAgentStat, 0, len(counts))

	for name, clicks := range counts {
		stats = append(stats, UserAgentStat{
			Name:   name,
			Clicks: clicks,
		})
	}

	sortUserAgentStats(stats)

	return stats, nil
}

func (r *ClickRepository) ListDevicesByShortCode(
	ctx context.Context,
	shortCode string,
) ([]UserAgentStat, error) {
	events, err := r.ListByShortCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf(
			"list devices: %w",
			err,
		)
	}

	counts := make(map[string]int64)

	for _, event := range events {
		device := DetectDevice(event.UserAgent)
		counts[device]++
	}

	stats := make([]UserAgentStat, 0, len(counts))

	for name, clicks := range counts {
		stats = append(stats, UserAgentStat{
			Name:   name,
			Clicks: clicks,
		})
	}

	sortUserAgentStats(stats)

	return stats, nil
}

func (r *ClickRepository) ListClientIPsByShortCode(
	ctx context.Context,
	shortCode string,
) ([]ClientIPStat, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			client_ip,
			COUNT(*) AS clicks
		FROM click_events
		WHERE short_code = ?
			AND TRIM(client_ip) <> ''
		GROUP BY client_ip
		ORDER BY clicks DESC, client_ip ASC
		`,
		shortCode,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list client IPs: %w",
			err,
		)
	}
	defer rows.Close()

	stats := make([]ClientIPStat, 0)

	for rows.Next() {
		var stat ClientIPStat

		if err := rows.Scan(
			&stat.IP,
			&stat.Clicks,
		); err != nil {
			return nil, fmt.Errorf(
				"scan client IP: %w",
				err,
			)
		}

		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate client IPs: %w",
			err,
		)
	}

	return stats, nil
}

func sortUserAgentStats(stats []UserAgentStat) {
	for i := 0; i < len(stats); i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[j].Clicks > stats[i].Clicks ||
				(stats[j].Clicks == stats[i].Clicks &&
					stats[j].Name < stats[i].Name) {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}
}

func (r *ClickRepository) ListByShortCode(
	ctx context.Context,
	shortCode string,
) ([]ClickEvent, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			id,
			short_code,
			clicked_at,
			referrer,
			user_agent,
			client_ip
		FROM click_events
		WHERE short_code = ?
		ORDER BY clicked_at DESC
		`,
		shortCode,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list click events: %w",
			err,
		)
	}
	defer rows.Close()

	events := make([]ClickEvent, 0)

	for rows.Next() {
		var event ClickEvent

		var clientIP sql.NullString

		if err := rows.Scan(
			&event.ID,
			&event.ShortCode,
			&event.ClickedAt,
			&event.Referrer,
			&event.UserAgent,
			&clientIP,
		); err != nil {
			return nil, fmt.Errorf(
				"scan click event: %w",
				err,
			)
		}

		if clientIP.Valid {
			event.ClientIP = clientIP.String
		}

		event.ClickedAt = event.ClickedAt.UTC()

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate click events: %w",
			err,
		)
	}

	return events, nil
}
