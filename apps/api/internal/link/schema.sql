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