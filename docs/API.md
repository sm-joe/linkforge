# LinkForge API

Base URL for local development:

```text
http://localhost:8080
```

## Health

### GET `/healthz`

Returns API health status.

## Links

### GET `/api/v1/links`

Returns available links.

### POST `/api/v1/links`

Creates a short link.

### GET `/api/v1/links/{id}`

Returns a link by ID.

### DELETE `/api/v1/links/{id}`

Deletes a link.

### POST `/api/v1/links/{id}/disable`

Disables a link.

### POST `/api/v1/links/{id}/enable`

Enables a link.

## Analytics

### GET `/api/v1/links/{id}/analytics`

Returns analytics for a link.

## QR Code

### GET `/api/v1/links/{id}/qr`

Returns a QR representation for a link.

## Health

### GET `/api/v1/links/{id}/health`

Checks destination health.

### GET `/api/v1/links/{id}/health/history`

Returns historical health information.

## Redirect

### GET `/{short_code}`

Redirects a short code to its destination when the link is active and has not expired.
