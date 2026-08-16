#!/usr/bin/env bash
# Add matches from data/cost.txt (year 2026 — weekday labels in cost.txt match 2026).
# Bills are verbatim from cost.txt; the free-match "cost lessened" amounts are
# recorded in notes only, NOT subtracted from total_bill.
# NOT idempotent — re-running creates duplicate matches.

set -euo pipefail

BASE_URL=${BASE_URL:-http://football.kubedb.cloud}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-changeme}

PLAYERS=$(curl -sS "$BASE_URL/api/players" | jq 'map({(.name): .id}) | add')

# add_match <date> <total_bill> <notes> <attendees>
# attendees: comma-separated names, optional ":<guest_count>" suffix (e.g. Mennal:1)
add_match() {
  local payload
  payload=$(jq -n --argjson map "$PLAYERS" \
    --arg date "$1" --argjson bill "$2" --arg notes "$3" --arg att "$4" '
    {
      date: $date,
      total_bill: $bill,
      notes: $notes,
      attendees: ($att | split(",") | map(
        split(":") as $p | {
          player_id: ($map[$p[0]] // error("unknown player: " + $p[0])),
          guest_count: (($p[1] // "0") | tonumber)
        }))
    }')
  curl -sS -X POST "$BASE_URL/api/matches" \
    -H "Authorization: Bearer $ADMIN_PASSWORD" \
    -H "Content-Type: application/json" \
    -d "$payload"
  echo
}

# ---- Card 1 ----
add_match 2026-06-04 3100 \
  'Free match (Card 1). Cost lessened elsewhere: 620 on Apr 16, 500 on May 14, 400 each on Jun 11, Jun 16, Jun 23, 600 on Jul 2' \
  'Arnob,Arman,Rabbi,Sadi,Biswarup,Rudro,Shohag,Amir'

add_match 2026-06-11 3340 \
  'Card 1. Free match (Jun 4) adjustment: 400 to be lessened' \
  'Shohag,Amir,Arman,Saurov,Rabbi,Sadi,Arnob,Biswarup,Imtiaz,Rasel,Rudro,Fakrul,Rubel,Zahin,Mamun'

add_match 2026-06-16 3280 \
  'Card 1. Free match (Jun 4) adjustment: 400 to be lessened' \
  'Arnob,Rabbi,Sadi,Arman,Biswarup,Rudro,Shohag,Amir,Hamim,Lotifur,Mennal'

add_match 2026-06-23 3320 \
  'Card 1. Free match (Jun 4) adjustment: 400 to be lessened' \
  'Arnob,Sadi,Saurov,Arman,Biswarup,Imtiaz,Amir,Shohag,Mennal'

add_match 2026-07-02 3080 \
  'Card 1. Free match (Jun 4) adjustment: 600 to be lessened' \
  'Arnob,Sadi,Saurov,Biswarup,Shohag,Mennal'

# ---- Card 2 ----
add_match 2026-07-08 3300 \
  'Card 2. Free match (Aug 4) adjustment: 400 to be lessened' \
  'Arnob,Rabbi,Saurov,Arman,Shohag,Mennal,Biswarup,Rudro'

add_match 2026-07-14 2825 \
  'Card 2. Free match (Aug 4) adjustment: 800 to be lessened' \
  'Arnob,Rabbi,Arman,Shohag,Mennal:1'

add_match 2026-07-23 3180 \
  'Card 2. Free match (Aug 4) adjustment: 400 to be lessened' \
  'Arnob,Rabbi,Arman,Shohag,Biswarup,Rasel,Rudro'

add_match 2026-07-28 3350 \
  'Card 2. Free match (Aug 4) adjustment: 400 to be lessened' \
  'Arnob,Rabbi,Arman,Shohag:6,Sadi'

add_match 2026-08-04 3900 \
  'Free match (Card 2). 3200 + 300 water already included in this bill. Cost lessened: 400 on Jul 8, 800 on Jul 14, 400 each on Jul 23, Jul 28, Aug 13' \
  'Avishek,Saber,Nirjhor,Shohag,Lotifur,Kawchar,Sabbir,Meraj,Sarwar,Shakil,Tarek,Arnob,Biswarup,Rabbi,Sadi'

add_match 2026-08-13 3300 \
  'Card 2. Free match (Aug 4) adjustment: 400 to be lessened' \
  'Arnob,Biswarup,Sadi,Arman,Shohag,Kawchar,Sarwar,Alfeh,Farzine,Meraj,Saber'
