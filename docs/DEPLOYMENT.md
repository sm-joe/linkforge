# LinkForge Deployment

## Requirements

- Docker or Podman
- Persistent storage for SQLite
- Network access between the web application and API

## API Container

Example:

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

## Web Container

The production web image runs Next.js using a dedicated non-root user.

The web application requires the API endpoint to be configured through:

```text
NEXT_PUBLIC_API_URL
```

## Persistence

The SQLite database is stored under:

```text
/app/data/linkforge.db
```

Mount persistent storage at `/app/data`.

Do not rely on container-local storage for production persistence.

## Health Checks

The API exposes:

```text
/healthz
```

The web container includes a health check against the local Next.js server.

## Production Considerations

For production deployments:

- Use persistent database storage.
- Place the application behind an appropriate reverse proxy or load balancer.
- Use HTTPS.
- Restrict network exposure where appropriate.
- Keep container images updated.
- Verify signed images before deployment.
- Back up the SQLite database.
- Monitor application and container health.
