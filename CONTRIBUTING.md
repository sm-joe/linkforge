# Contributing to LinkForge

## Development Workflow

1. Fork or clone the repository.
2. Create a focused branch for your change.
3. Make the smallest appropriate change.
4. Run the relevant tests and validation.
5. Review the final diff.
6. Open a pull request with a clear description.

## Backend

From `apps/api`:

```powershell
go test ./...
go vet ./...
```

## URL Validation Package

From `packages/url-validation`:

```powershell
go test ./...
go vet ./...
```

## Frontend

From `apps/web`:

```powershell
npm install
npm run lint
npm run build
```

## Security

Changes must not weaken existing security controls without a documented reason.

Relevant security checks include:

- Secret scanning
- Static analysis
- Dependency scanning
- IaC scanning
- Container scanning
- Image signing and verification

## Pull Requests

Pull requests should:

- Explain what changed
- Explain why it changed
- Keep unrelated changes out
- Include relevant test results
- Include documentation updates when behavior changes

## Commits

Use clear, concise commit messages that describe the change.

Avoid committing:

- Secrets
- Credentials
- Local databases
- Build artifacts
- Temporary files
- Environment-specific configuration containing sensitive values
