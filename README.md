# LinkForge

LinkForge is a portable, self-hosted URL shortener built with Go, Next.js, and SQLite.

It provides a lightweight interface for creating and managing short links while keeping the application easy to run locally or deploy in containers.

## Features

- Shorten long URLs
- Custom aliases
- Destination URL validation
- Safe destination validation against private/reserved IP addresses
- SSRF protection during link health checks
- Link expiration
- Enable and disable links
- Link deletion
- Link health checks
- Link health history
- Click analytics
- Client IP tracking
- Referrer tracking
- User-agent tracking
- Browser analytics
- QR code generation
- REST API
- SQLite persistence
- Docker/Podman deployment
- Container hardening
- API rate limiting
- Security response headers
- CORS protection
- Automated security scanning in CI/CD

## Architecture

```text
Browser
   |
   v
Next.js Web
   |
   v
Go API
   |
   +--> SQLite
   +--> URL Validation
   +--> Health / SSRF Validation
```

## Project Structure

```text
linkforge/
├── apps/
│   ├── api/
│   └── web/
├── packages/
│   └── url-validation/
├── deploy/
├── deployment/
├── .github/
│   └── workflows/
├── docs/
├── README.md
├── SECURITY.md
├── CONTRIBUTING.md
├── CHANGELOG.md
└── LICENSE
```

## Technology Stack

### Backend
- Go
- SQLite
- REST API

### Frontend
- Next.js
- React
- TypeScript

### Infrastructure
- Docker / Podman
- Docker Compose
- GitHub Actions

### Security
- Gitleaks
- CodeQL
- OSV Scanner
- Checkov
- Trivy
- Cosign

## Prerequisites

Install:

- Go
- Node.js
- npm
- Docker or Podman
- Git

## Local Development

### Backend

```powershell
cd apps\api
go test ./...
go vet ./...
```

The API listens on:

```text
http://localhost:8080
```

### Frontend

```powershell
cd apps\web
npm install
npm run dev
```

The web application is available at:

```text
http://localhost:3000
```

The frontend uses `NEXT_PUBLIC_API_URL` to determine the API endpoint.

## API

```text
GET    /healthz

GET    /api/v1/links
POST   /api/v1/links

GET    /api/v1/links/{id}
DELETE /api/v1/links/{id}

POST   /api/v1/links/{id}/disable
POST   /api/v1/links/{id}/enable

GET    /api/v1/links/{id}/analytics
GET    /api/v1/links/{id}/qr
GET    /api/v1/links/{id}/health
GET    /api/v1/links/{id}/health/history

GET    /{short_code}
```

See `docs/API.md` for details.

## URL Validation and SSRF Protection

Destination URLs are validated before storage.

Health checks additionally validate resolved destinations, block private/reserved addresses, disable environment proxies, require TLS 1.2 or newer, and do not follow redirects.

## Rate Limiting

API requests are rate limited to reduce abuse and excessive request volume. Health endpoints are exempt for monitoring.

## Container Security

Production web containers run as a dedicated non-root user. Container images are scanned as part of CI/CD.

## Persistence

LinkForge uses SQLite. Container deployments should use persistent storage for the database.

## Docker / Podman

Example API image build:

```powershell
podman build -f deploy\docker\api\Dockerfile -t localhost/linkforge-api:latest .
```

Example API container:

```powershell
podman run -d `
  --name linkforge-api `
  --cgroups=disabled `
  -p 8080:8080 `
  -v linkforge-data:/app/data `
  -e PORT=8080 `
  -e DATABASE_PATH=/app/data/linkforge.db `
  localhost/linkforge-api:latest
```

See `docs/DEPLOYMENT.md` for deployment details.

## Testing

```powershell
cd apps\api
go test ./...
go vet ./...
```

```powershell
cd packages\url-validation
go test ./...
go vet ./...
```

```powershell
cd apps\web
npm run lint
npm run build
```

## CI/CD

GitHub Actions performs automated validation and security checks including:

- Go tests
- Go vet
- Secret scanning
- Code analysis
- Dependency scanning
- IaC scanning
- Filesystem scanning
- Container image scanning
- Image signing
- Image verification
- Image publishing

Workflow definitions are under `.github/workflows/`.

## Security

See `SECURITY.md` for the security policy and vulnerability reporting process.

## Contributing

See `CONTRIBUTING.md`.

## License

LinkForge is licensed under the Apache License 2.0.
