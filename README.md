# football_calc

Go HTTP server for tracking weekly turf football match costs, dues, payments, and balances.

## Environment variables

Loaded from a `.env` file at startup (via `godotenv`) or from the process environment.

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | yes | — | PostgreSQL connection string, e.g. `postgres://user:pass@localhost:5432/football_calc?sslmode=disable` |
| `ADMIN_PASSWORD` | yes | — | Bearer token required for all write endpoints and `/api/admin/*` |
| `PORT` | no | `8080` | HTTP listen port |

## Local run

```bash
cp .env.example .env   # then edit DATABASE_URL and ADMIN_PASSWORD
make run               # go run ./cmd/server serve
```

The server bootstraps the schema on startup and listens on `PORT` (default 8080).

## Common commands

```bash
make build      # build static binary to bin/football-calc
make run        # run locally (reads .env)
make fmt        # gofmt + goimports
make lint       # golangci-lint
make container  # docker buildx build
make push       # build + push image to ghcr.io
```
