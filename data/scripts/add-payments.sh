#!/usr/bin/env bash
# Add payments from the "Card 3" block of data/cost.txt (year 2026).
# Recorded as standalone payments (no match_id) — cost.txt does not tie them to a match.
# NOT idempotent — re-running creates duplicate payments.

set -euo pipefail

BASE_URL=${BASE_URL:-http://football.kubedb.cloud}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-changeme}

PLAYERS=$(curl -sS "$BASE_URL/api/players" | jq 'map({(.name): .id}) | add')

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
  curl -sS -X POST "$BASE_URL/api/payments" \
    -H "Authorization: Bearer $ADMIN_PASSWORD" \
    -H "Content-Type: application/json" \
    -d "$payload"
  echo
}

# Listed under the "Card 3" heading in cost.txt, but noted by the card each date falls in.
add_payment Shahed  219  2026-06-11 'Card 1'
add_payment Palash  281  2026-06-11 'Card 1'
add_payment Shohag  3000 2026-06-12 'Card 1'
add_payment Mennal  400  2026-06-23 'Card 1'
add_payment Shohag  1000 2026-07-08 'Card 2'
add_payment Saurov  1000 2026-07-08 'Card 2'
add_payment Rabbi   1000 2026-07-08 'Card 2'

# cost.txt lists "Amir 2000?" — amount unconfirmed and no date. Uncomment once confirmed.
# add_payment Amir 2000 <paid_at> '<card>'
