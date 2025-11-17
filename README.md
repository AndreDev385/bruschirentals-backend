# Bruschi Rentals Backend

A Go-based REST API backend for managing client lifecycles and rentals in a rental company. Built with Echo, PostgreSQL, OpenTelemetry tracing, and Zap logging.

## Features

- Health check endpoint
- Database migrations with golang-migrate
- Structured logging with Zap
- Distributed tracing with OpenTelemetry (OTLP)
- Docker support for dev and prod
- Hot reloading with Air

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL (for local dev)

## Quick Start

1. **Clone and setup:**
   ```bash
   git clone <repo-url>
   cd bruschi-rentals-backend
   cp .env.example .env
   # Edit .env with your DATABASE_URL
   ```

2. **Run with Docker (recommended):**
   ```bash
   make docker-dev  # For development with hot reload
   # or
   make docker-prod  # For production build
   ```

3. **Run locally:**
   ```bash
   # Install dependencies
   go mod tidy

   # Run migrations
   make migrate-up

   # Start server
   make dev  # With hot reload
   # or
   make run  # Without
   ```

4. **Test:**
   ```bash
   make test
   ```

## API Endpoints

- `GET /api/v1/health` - Health check (returns DB status)
- `GET /swagger/*` - API documentation (Swagger UI)

### Example Requests

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Response: {"status": "healthy"}
```

View full API docs at `http://localhost:8080/swagger/index.html`.

## Configuration

Environment variables (see `.env.example`):

- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - Server port (default: 8080)
- `OTEL_EXPORTER_OTLP_ENDPOINT` - OTLP endpoint for tracing (optional)
- `ENV` - Environment (development/production)

## Database

Migrations are in `/migrations`. Run with:
```bash
make migrate-up    # Apply migrations
make migrate-down  # Rollback
make migrate-create NAME=migration_name  # Create new migration
```

## Development

- Use `make dev` for hot reloading.
- Tracing: View in Jaeger at http://localhost:16686 (dev only).
- Logs: Structured JSON to stdout.

## Deployment

### Railway (Recommended)

1. Connect GitHub repo to Railway.
2. Add Postgres service in Railway dashboard.
3. Set env vars: `DATABASE_URL` (from Postgres service), `OTEL_EXPORTER_OTLP_ENDPOINT` (optional for tracing).
4. Railway auto-builds from source and runs migrations on startup.
5. Deploy automatically on push to main.

### Local Production

Build and run:
```bash
go build -o server ./cmd/server
migrate -path ./migrations -database "$DATABASE_URL" up
./server
```

## Project Structure

```
.
├── cmd/server/          # Main application entry point
├── internal/            # Private application code
│   ├── config/          # Configuration loading
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Custom middlewares
│   ├── logging/         # Logging setup
│   └── tracing/         # Tracing setup
├── migrations/          # DB migrations
├── pkg/                 # (Future: public libraries)
├── Dockerfile.dev       # Dev Docker image
├── docker-compose.dev.yml # Dev compose
├── railway.toml         # Railway config
├── Makefile             # Build tasks
├── .env.example         # Env template
└── README.md            # This file
```

## Troubleshooting

- **DB connection fails:** Check `DATABASE_URL` in `.env`.
- **Migrations fail:** Ensure DB is running and accessible.
- **Tracing not working:** Verify `OTEL_EXPORTER_OTLP_ENDPOINT` is set.
- **Build fails:** Run `go mod tidy` and check Go version (1.21+).

## Contributing

- Follow Go best practices.
- Run tests and lint before PRs.
- Use conventional commits.

## License

MIT
