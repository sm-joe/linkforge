package link

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func OpenDatabase(path string) (*sql.DB, error) {
	if path == "" {
		path = "./data/linkforge.db"
	}

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf(
			"create database directory: %w",
			err,
		)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf(
			"open database: %w",
			err,
		)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf(
			"ping database: %w",
			err,
		)
	}

	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf(
			"initialize schema: %w",
			err,
		)
	}

	return db, nil
}

const schema = `
CREATE TABLE IF NOT EXISTS links (
    id TEXT PRIMARY KEY,
    short_code TEXT NOT NULL UNIQUE,
    alias TEXT,
    destination TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    expires_at DATETIME,
		client_ip TEXT
);

CREATE INDEX IF NOT EXISTS idx_links_short_code
    ON links(short_code);

CREATE INDEX IF NOT EXISTS idx_links_status
    ON links(status);

CREATE TABLE IF NOT EXISTS click_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_code TEXT NOT NULL,
    clicked_at DATETIME NOT NULL,
    referrer TEXT,
    user_agent TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_click_events_short_code
    ON click_events(short_code);

CREATE INDEX IF NOT EXISTS idx_click_events_clicked_at
    ON click_events(clicked_at);

CREATE TABLE IF NOT EXISTS link_health_checks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_code TEXT NOT NULL,
    destination TEXT NOT NULL,
    status TEXT NOT NULL,
    http_status INTEGER,
    response_time_ms INTEGER NOT NULL,
    final_url TEXT,
    redirects INTEGER NOT NULL,
    https BOOLEAN NOT NULL,
    checked_at DATETIME NOT NULL,
    error TEXT
);

CREATE INDEX IF NOT EXISTS idx_link_health_checks_short_code
    ON link_health_checks(short_code);

CREATE INDEX IF NOT EXISTS idx_link_health_checks_checked_at
    ON link_health_checks(checked_at);
`
