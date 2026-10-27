#!/usr/bin/env bash
# Add payments from data/new-entries.txt (year 2026).
# Recorded as standalone payments (no match_id), dated on the match they were listed under.
# NOT idempotent — re-running creates duplicate payments.

set -euo pipefail

BASE_URL=${BASE_URL:-http://football.kubedb.cloud}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-changeme}

PLAYERS=$(curl -sS --fail-with-body "$BASE_URL/api/players" | jq 'map({(.name): .id}) | add')

# add_payment <player> <amount> <paid_at> [notes]
add_payment() {
  local payload
  payload=$(jq -n --argjson map "$PLAYERS" \
    --arg name "$1" --argjson amount "$2" --arg paid_at "$3" --arg notes "${4:-}" '
    {
      player_id: ($map[$name] // error("unknown player: " + $name)),
      amount: $amount,
      paid_at: $paid_at,
      notes: (if $notes == "" then null else $notes end)
    }')
  curl -sS --fail-with-body -X POST "$BASE_URL/api/payments" \
    -H "Authorization: Bearer $ADMIN_PASSWORD" \
    -H "Content-Type: application/json" \
    -d "$payload"
  echo
}

add_payment Nirjhor 260 2026-08-16 'bKash'
add_payment Arman  1000 2026-08-20
add_payment Rizu   1900 2026-08-20 'Paid 2000; took 100 back on Aug 27'
add_payment Sarwar 500  2026-08-20
add_payment Alfeh  1000 2026-08-20
add_payment Tarek  260  2026-08-20
add_payment Shakil 260  2026-08-20
add_payment Hamim  583  2026-08-25
add_payment Rabbi  500  2026-09-02 'Advance'
add_payment Rabbi  500  2026-09-08 'Advance'
add_payment Shohag 500  2026-09-12 'Advance'
add_payment Shohag 500  2026-09-26 'Advance'
