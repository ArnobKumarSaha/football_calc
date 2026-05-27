# Plan: Football Calc — Go Web App

## Context
Weekly turf football matches among a fixed group. Admin records each match's total bill, attendees (with optional guest counts), and payments. App computes per-person dues, tracks running balances (credit/debt), and exposes reports. Calculation starts fresh — no historical data, no migration tooling. Schema bootstrapped via `CREATE TABLE IF NOT EXISTS` on startup. Deployed on K8s with an **external managed PostgreSQL**, exposed via Gateway API + LoadBalancer.

---

## Data Model (PostgreSQL)

```sql
players (
  id          SERIAL PRIMARY KEY,
  name        TEXT NOT NULL UNIQUE,
  created_at  TIMESTAMPTZ DEFAULT now(),
  deleted_at  TIMESTAMPTZ
)

matches (
  id          SERIAL PRIMARY KEY,
  date        DATE NOT NULL,
  total_bill  NUMERIC(10,2) NOT NULL,
  notes       TEXT,
  created_at  TIMESTAMPTZ DEFAULT now(),
  deleted_at  TIMESTAMPTZ
)

match_attendees (
  id           SERIAL PRIMARY KEY,
  match_id     INT NOT NULL REFERENCES matches(id),
  player_id    INT NOT NULL REFERENCES players(id),
  guest_count  INT NOT NULL DEFAULT 0,
  UNIQUE(match_id, player_id)
)

payments (
  id          SERIAL PRIMARY KEY,
  player_id   INT NOT NULL REFERENCES players(id),
  match_id    INT REFERENCES matches(id),   -- NULL = standalone settlement
  amount      NUMERIC(10,2) NOT NULL,
  notes       TEXT,
  paid_at     DATE NOT NULL,
  created_at  TIMESTAMPTZ DEFAULT now()
)
```

Schema bootstrap: single `schema.sql` executed on startup. No migration library.

---

## Balance Logic

Per match:
- `units = SUM(1 + guest_count)` over attendees
- `per_unit = total_bill / units`
- Each attendee `raw_due = per_unit * (1 + guest_count)`
- **Rounding**: each `raw_due` rounded to 2 decimals; the **last attendee** (by `match_attendees.id`) absorbs the remainder so `SUM(due) == total_bill` exactly.

Running balance per player (ignoring soft-deleted rows):
```
balance = SUM(payments.amount) - SUM(dues)
```
- Positive = credit (overpaid)
- Negative = debt (owes)
- Standalone payments (NULL `match_id`) contribute to `SUM(payments.amount)` only.
- **New players** (not in any match or payment): starting balance = **0**

---

## Mutability
- Matches, attendees, payments, players support full **edit + delete**.
- Delete is **soft** (`deleted_at IS NOT NULL`); calc queries filter these out.

---

## Project Structure

```
football-calc/
├── cmd/server/main.go      # cobra root cmd, starts HTTP server
├── pkg/
│   ├── api/                # chi handlers (one file per resource)
│   ├── db/                 # pgxpool + queries + tx helpers + schema bootstrap
│   ├── models/             # domain structs
│   └── service/            # balance calc, rounding, audit wrapper
├── schema.sql              # CREATE TABLE IF NOT EXISTS ...
├── static/index.html       # vanilla JS SPA (embedded via go:embed)
├── deploy/
│   ├── namespace.yaml
│   ├── secret.yaml         # ADMIN_PASSWORD, DATABASE_URL
│   ├── deployment.yaml     # probes + resource limits
│   ├── service.yaml
│   └── httproute.yaml
├── Dockerfile
├── Makefile                # build / lint / fmt / container / push
├── go.mod
└── vendor/
```

---

## API

All non-GET endpoints require `Authorization: Bearer <ADMIN_PASSWORD>`.

### Players
| Method | Path | Body / Notes |
|--------|------|--------------|
| GET    | `/api/players` | list + current balance |
| POST   | `/api/players` | `{name}` |
| PATCH  | `/api/players/:id` | `{name}` |
| DELETE | `/api/players/:id` | soft-delete |
| GET    | `/api/players/:id/balance` | full per-match history |

### Matches
| Method | Path | Body / Notes |
|--------|------|--------------|
| GET    | `/api/matches` | list summaries (paginated) |
| POST   | `/api/matches` | `{date, total_bill, notes}` |
| GET    | `/api/matches/:id` | detail: attendees + dues + payments |
| PATCH  | `/api/matches/:id` | `{date?, total_bill?, notes?}` |
| DELETE | `/api/matches/:id` | soft-delete |

### Attendees
| Method | Path | Body / Notes |
|--------|------|--------------|
| POST   | `/api/matches/:id/attendees` | `{player_id, guest_count}` (upsert) |
| PATCH  | `/api/matches/:id/attendees/:player_id` | `{guest_count}` |
| DELETE | `/api/matches/:id/attendees/:player_id` | |

### Payments
| Method | Path | Body / Notes |
|--------|------|--------------|
| GET    | `/api/payments` | filter `?player_id=&match_id=` |
| POST   | `/api/payments` | `{player_id, amount, match_id?, paid_at, notes?}` |
| PATCH  | `/api/payments/:id` | any field |
| DELETE | `/api/payments/:id` | |

### Reports
| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/reports/monthly` | `?year=2026` → total spent, match count per month |
| GET | `/api/reports/attendance` | per-player attendance rate |
| GET | `/api/reports/matches.csv` | CSV: matches + attendees + dues |
| GET | `/api/reports/payments.csv` | CSV: all payments |

### Ops
| Method | Path | Notes |
|--------|------|-------|
| GET | `/healthz` | liveness |
| GET | `/readyz`  | DB ping |

---

## Tech Stack

| Concern | Choice |
|---------|--------|
| Router | `github.com/go-chi/chi/v5` |
| DB | `github.com/jackc/pgx/v5` (pgxpool) |
| Schema bootstrap | raw `schema.sql` on startup |
| Config | env: `DATABASE_URL`, `ADMIN_PASSWORD`, `PORT` |
| Logging | `k8s.io/klog/v2` (structured) |
| CLI | `github.com/spf13/cobra` |
| Frontend | single `index.html` + vanilla JS, `go:embed` |

---

## Auth
Middleware on all non-GET routes checks `Authorization: Bearer <token>` against `ADMIN_PASSWORD`. Returns `401` otherwise.

---

## Frontend (Vanilla JS SPA)

Sections:
1. **Matches** — list + expand for breakdown
2. **Record match** — date, bill, notes
3. **Attendance** — pick players + guest counts for a match
4. **Payments** — record per-match or standalone payment (match dropdown optional)
5. **Balances** — table: name, total due, total paid, balance
6. **Reports** — monthly summary, attendance rate, CSV download links

Admin password held in `sessionStorage`; sent as Bearer on writes.

---

## K8s Deployment

- Namespace: `football-calc`
- Secret `football-calc-secret`: `ADMIN_PASSWORD`, `DATABASE_URL` (→ **external managed Postgres**)
- Deployment: 1 replica, env from secret, liveness `/healthz`, readiness `/readyz`, resource requests/limits set
- Service: ClusterIP :8080
- HTTPRoute (`gateway.networking.k8s.io/v1`) → existing LB gateway

---

## Verification

1. `go build ./cmd/server` — compiles clean
2. Local run: `DATABASE_URL=... ADMIN_PASSWORD=test go run ./cmd/server` — verify schema auto-created, `/healthz` + `/readyz` OK
3. Smoke (curl):
   - Create 4 players → create match (bill=1000) → add 4 attendees (one with guest_count=1) → verify `units=5`, `per_unit=200`, dues `[200,200,200,400]`, **sum == 1000**
   - Record 2 partial payments for one player on a match → verify balance
   - Record a standalone payment (no match) → verify it shifts balance only
   - Edit `total_bill` → balances recalc
4. Reports: `/api/reports/monthly?year=2026`, `/api/reports/attendance`, `/api/reports/matches.csv` open in browser/Excel
5. `kubectl apply -k deploy/` → HTTPRoute reachable via LB; pod healthy
