<p align="center">
  <strong>Portable, self-hosted URL shortener for simple, secure link management.</strong><br>
  Create short links with custom aliases, expiration, analytics, QR codes,
  destination validation, and SSRF-resistant health checks.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.27-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.27">
  <img src="https://img.shields.io/badge/Next.js-16-000000?style=for-the-badge&logo=next.js&logoColor=white" alt="Next.js 16">
  <img src="https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react&logoColor=black" alt="React 19">
  <img src="https://img.shields.io/badge/TypeScript-5-3178C6?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript">
  <img src="https://img.shields.io/badge/SQLite-3-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Self--Hosted-2ea44f?style=flat-square" alt="Self Hosted">
  <img src="https://img.shields.io/badge/SSRF-Protected-2088FF?style=flat-square" alt="SSRF Protected">
  <img src="https://img.shields.io/badge/Rate--Limited-6f42c1?style=flat-square" alt="Rate Limited">
  <img src="https://img.shields.io/badge/Container-Hardened-1904DA?style=flat-square" alt="Container Hardened">
  <img src="https://img.shields.io/badge/CI%2FCD-Security%20Gated-success?style=flat-square" alt="CI/CD Security Gated">
</p>

# LinkForge

LinkForge is a portable, self-hosted URL shortener built with Go, Next.js, and SQLite.

It provides a lightweight interface and REST API for creating and managing short links while keeping deployment simple, portable, and security-focused.

## Features

- Shorten long URLs
- Custom aliases
- Link expiration
- Enable and disable links
- Link deletion
- Destination URL validation
- Private and reserved IP protection
- SSRF-resistant link health checks
- Link health history
- Click analytics
- Client IP tracking
- Referrer and browser analytics
- QR code generation
- REST API
- SQLite persistence
- Docker / Podman deployment
- Non-root container execution
- API rate limiting
- Security response headers
- CORS protection
- Automated CI/CD security gates
- Container image scanning
- Container image signing and verification

## Architecture

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
                  +----------+        +-----------+
                  |                               |
                  v                               v
            +-----------+                  +---------------+
            |   SQLite  |                  | URL Validation|
            +-----------+                  +-------+-------+
                                                    |
                                                    v
                                             +-------------+
                                             | SSRF / DNS  |
                                             | Protection  |
                                             +-------------+
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
├── docs/
├── .github/
│   └── workflows/
├── .env.example
├── CHANGELOG.md
├── CONTRIBUTING.md
├── LICENSE
├── README.md
└── SECURITY.md
```

## Technology Stack

| Component | Technology |
|---|---|
| Backend | Go |
| Frontend | Next.js / React / TypeScript |
| Database | SQLite |
| Containers | Docker / Podman |
| CI/CD | GitHub Actions |
| Secret Scanning | Gitleaks |
| Code Analysis | CodeQL |
| Dependency Scanning | OSV Scanner |
| IaC Scanning | Checkov |
| Container Scanning | Trivy |
| Image Signing | Cosign |

## Security

LinkForge treats user-supplied destination URLs as untrusted input.

Security controls include:

- URL scheme validation
- Hostname validation
- DNS validation
- Private IP blocking
- Reserved IP blocking
- DNS rebinding protection during health checks
- Redirect blocking during health checks
- Environment proxy disabled for health requests
- TLS 1.2 minimum
- API rate limiting
- Security response headers
- Non-root container execution
- Automated security scanning
- Container image signing and verification

See [`SECURITY.md`](SECURITY.md) for the security policy.

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

See [`docs/API.md`](docs/API.md) for the API reference.

## Local Development

### Backend

```powershell
cd apps\api
go test ./...
go vet ./...
```

### Frontend

```powershell
cd apps\web
npm install
npm run dev
```

The web application runs on:

```text
http://localhost:3000
```

The API runs on:

```text
http://localhost:8080
```

The frontend API endpoint is configured through:

```text
NEXT_PUBLIC_API_URL
```

See [`.env.example`](.env.example) for the available configuration.

## Container Deployment

LinkForge supports Docker and Podman-based deployment.

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

See [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md) for deployment details.

## Testing

Backend:

```powershell
cd apps\api
go test ./...
go vet ./...
```

URL validation:

```powershell
cd packages\url-validation
go test ./...
go vet ./...
```

Frontend:

```powershell
cd apps\web
npm run lint
npm run build
```

## CI/CD

GitHub Actions provides automated:

- Testing
- Static analysis
- Secret scanning
- Dependency scanning
- IaC scanning
- Container scanning
- Image building
- Image publishing
- Image signing
- Image verification

Workflow definitions are located under `.github/workflows/`.

## Documentation

- [`SECURITY.md`](SECURITY.md) — security policy
- [`CONTRIBUTING.md`](CONTRIBUTING.md) — contribution guidelines
- [`CHANGELOG.md`](CHANGELOG.md) — project changes
- [`docs/API.md`](docs/API.md) — API reference
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — architecture
- [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md) — deployment guide

## License

LinkForge is licensed under the Apache License 2.0.
