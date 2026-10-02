# Security Policy

## Scope

Security issues affecting LinkForge, its API, web application, container deployment, or supporting infrastructure are in scope.

## Security Controls

### URL Validation

Destination URLs are validated before acceptance, including hostname, scheme, IP, private-address, reserved-address, and DNS checks.

### SSRF Protection

Health checks:

- Disable environment proxies
- Require TLS 1.2 or newer
- Validate resolved addresses
- Block private/reserved destinations
- Revalidate destinations during connection
- Do not follow redirects

### Rate Limiting

API requests are rate limited to reduce abuse and excessive request volume.

### Security Headers

The API provides security headers including:

- `X-Content-Type-Options`
- `X-Frame-Options`
- `Referrer-Policy`
- `Permissions-Policy`

### Container Security

Production containers use a dedicated non-root runtime user. Container images are scanned through CI/CD.

## Security Testing

CI/CD includes:

- Secret scanning
- Static analysis
- Dependency scanning
- IaC scanning
- Container filesystem scanning
- Container image vulnerability scanning
- Image signing
- Image verification

## Reporting a Vulnerability

Please report suspected security vulnerabilities privately rather than opening a public issue.

Include:

- Description
- Affected component
- Reproduction steps
- Potential impact
- Relevant safe-to-share evidence

Do not include credentials, tokens, personal information, or other sensitive data.

## Supported Versions

Security fixes are applied to actively maintained versions of LinkForge. Upgrade older deployments to the current maintained release before reporting issues that may already be fixed.
