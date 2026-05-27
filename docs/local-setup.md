# Local Setup

## Prerequisites

- Go 1.24+
- PostgreSQL (local or remote)

## Run

```bash
# Set required env vars
export DATABASE_URL="postgres://user:password@localhost:5432/football_calc?sslmode=disable"
export ADMIN_PASSWORD="changeme"
export PORT=9093

# Build and run
make build
./bin/football-calc

# Or directly
go run -mod=vendor ./cmd/server
```

Server starts on `:8080` by default. Override with `PORT=<n>`.

Schema is auto-bootstrapped on first start (`CREATE TABLE IF NOT EXISTS`).

## Verify

```bash
# Health checks
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz

# Create a player
curl -X POST http://localhost:8080/api/players \
  -H "Authorization: Bearer changeme" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice"}'

# Create a match
curl -X POST http://localhost:8080/api/matches \
  -H "Authorization: Bearer changeme" \
  -H "Content-Type: application/json" \
  -d '{"date":"2026-05-27","total_bill":1000,"notes":"turf"}'

# Add attendees (4 players, one with guest_count=1 → 5 units, per_unit=200)
curl -X POST http://localhost:8080/api/matches/1/attendees \
  -H "Authorization: Bearer changeme" \
  -H "Content-Type: application/json" \
  -d '{"player_id":1,"guest_count":0}'

# Get balances
curl http://localhost:8080/api/players
```

## Frontend

Open `http://localhost:8080` in a browser. Admin password is stored in `sessionStorage`.
