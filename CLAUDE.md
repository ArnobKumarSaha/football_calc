# football-calc

Go HTTP server for tracking weekly turf football match costs, dues, payments, and balances.

## What it does

- Admin records matches (date, total bill), attendees (with optional guest counts), and payments
- Computes per-person dues with rounding (last attendee absorbs remainder so sum == bill exactly)
- Running balance per player: `SUM(payments) - SUM(dues)` — positive = credit, negative = debt
- Standalone payments (no match) contribute to payments only
- All deletes are soft (`deleted_at IS NOT NULL`); calc queries filter these out

## Tech Stack

| Concern | Choice |
|---------|--------|
| Router | `github.com/go-chi/chi/v5` |
| DB | `github.com/jackc/pgx/v5` (pgxpool) |
| Schema bootstrap | `pkg/db/schema.sql` run on startup via `db.Bootstrap` |
| Config | env vars: `DATABASE_URL`, `ADMIN_PASSWORD`, `PORT` (default 8080) |
| Logging | `k8s.io/klog/v2` |
| CLI | `github.com/spf13/cobra` |
| Frontend | `cmd/server/static/index.html` (vanilla JS SPA, `go:embed`) |

## Project Layout

```
cmd/server/main.go          # cobra root, connects DB, starts HTTP
cmd/server/static/          # embedded frontend (go:embed)
pkg/api/                    # chi handlers, one file per resource
pkg/db/                     # pgxpool, schema bootstrap, queries
pkg/models/                 # domain structs
pkg/service/dues.go         # balance calc + rounding logic
deploy/                     # K8s manifests
schema.sql                  # top-level copy (canonical is pkg/db/schema.sql)
```

## Auth

All non-GET endpoints require `Authorization: Bearer <ADMIN_PASSWORD>`. Returns 401 otherwise.

## API Surface

### Players
| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/players` | list + current balance |
| POST | `/api/players` | `{name}` |
| PATCH | `/api/players/:id` | `{name}` |
| DELETE | `/api/players/:id` | soft-delete |
| GET | `/api/players/:id/balance` | full per-match history |

### Matches
| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/matches` | paginated list |
| POST | `/api/matches` | `{date, total_bill, notes}` |
| GET | `/api/matches/:id` | detail: attendees + dues + payments |
| PATCH | `/api/matches/:id` | `{date?, total_bill?, notes?}` |
| DELETE | `/api/matches/:id` | soft-delete |

### Attendees
| Method | Path | Notes |
|--------|------|-------|
| POST | `/api/matches/:id/attendees` | `{player_id, guest_count}` upsert |
| PATCH | `/api/matches/:id/attendees/:player_id` | `{guest_count}` |
| DELETE | `/api/matches/:id/attendees/:player_id` | |

### Payments
| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/payments` | filter `?player_id=&match_id=` |
| POST | `/api/payments` | `{player_id, amount, match_id?, paid_at, notes?}` |
| PATCH | `/api/payments/:id` | any field |
| DELETE | `/api/payments/:id` | |

### Reports
| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/reports/monthly` | `?year=YYYY` |
| GET | `/api/reports/attendance` | per-player attendance rate |
| GET | `/api/reports/matches.csv` | CSV export |
| GET | `/api/reports/payments.csv` | CSV export |

### Ops
| Method | Path |
|--------|------|
| GET | `/healthz` |
| GET | `/readyz` |

## Common Commands

```bash
make build      # CGO_ENABLED=0 go build -mod=vendor -o bin/football-calc ./cmd/server
make fmt        # gofmt + goimports
make lint       # golangci-lint
make container  # docker buildx build (linux/amd64,linux/arm64)
make push       # build + push to ghcr.io/arnobkumarsaha/football-calc:latest
```

## K8s Deployment

- Namespace: `football-calc`
- Secret `football-calc-secret`: keys `ADMIN_PASSWORD`, `DATABASE_URL`
- External managed PostgreSQL (not in-cluster)
- Exposed via Gateway API `HTTPRoute` → `lb-gateway` in namespace `gateway`
- Liveness: `/healthz`, Readiness: `/readyz`
- Resources: 50m/64Mi requests, 200m/128Mi limits
