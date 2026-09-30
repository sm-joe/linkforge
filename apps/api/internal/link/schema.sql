CREATE TABLE IF NOT EXISTS links (
    id TEXT PRIMARY KEY,
    short_code TEXT NOT NULL UNIQUE,
    alias TEXT,
    destination TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    expires_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_links_short_code
    ON links(short_code);

CREATE INDEX IF NOT EXISTS idx_links_status
    ON links(status);