# LinkForge Architecture

## Overview

LinkForge consists of a Next.js frontend, Go API, SQLite persistence, and a shared URL-validation package.

```text
                    +----------------+
                    |    Browser     |
                    +-------+--------+
                            |
                            v
                    +----------------+
                    |   Next.js Web  |
                    +-------+--------+
                            |
                            v
                    +----------------+
                    |     Go API     |
                    +---+--------+---+
                        |        |
                        v        v
                   +--------+ +-----------+
                   | SQLite | | URL / SSRF |
                   |        | | Validation |
                   +--------+ +-----------+
```

## Components

### Web

`apps/web`

Provides the user interface for creating, viewing, and managing links.

### API

`apps/api`

Provides the REST API, redirect handling, link management, analytics, QR generation, and health checks.

### URL Validation

`packages/url-validation`

Provides reusable destination validation and private/reserved address protection.

### Deployment

`deploy/` and `deployment/`

Contain container and deployment-related configuration and automation.

### CI/CD

`.github/workflows/`

Contains automated build, test, security scanning, image signing, verification, and publishing workflows.

## Data

SQLite stores application state and analytics data.

Container deployments should mount persistent storage for the database.

## Security Boundary

User-supplied destination URLs are treated as untrusted input.

URL validation occurs before persistence, while health checks apply additional network-level destination validation to reduce SSRF risk.
