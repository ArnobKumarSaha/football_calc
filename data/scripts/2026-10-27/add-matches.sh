#!/usr/bin/env bash
# Add matches from data/new-entries.txt (year 2026).
# Name mapping: Sohag -> Shohag, Kausar -> Kawchar, Asad / "Asad bin saber" -> Saber.
# NOT idempotent — re-running creates duplicate matches.

set -euo pipefail

BASE_URL=${BASE_URL:-http://football.kubedb.cloud}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-changeme}

PLAYERS=$(curl -sS --fail-with-body "$BASE_URL/api/players" | jq 'map({(.name): .id}) | add')

# add_match <date> <total_bill> <notes> <attendees>
# attendees: comma-separated names, optional ":<guest_count>" suffix (e.g. Rizu:1)
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
  curl -sS --fail-with-body -X POST "$BASE_URL/api/matches" \
    -H "Authorization: Bearer $ADMIN_PASSWORD" \
    -H "Content-Type: application/json" \
    -d "$payload"
  echo
}

add_match 2026-08-21 4300 '' \
  'Arnob,Rabbi,Arman,Sadi,Rasel,Biswarup:1,Shohag,Rizu,Sarwar,Alfeh,Naimun,Saimun,Fahim'

add_match 2026-08-27 4120 \
  'Bill 4000 + 120. Rizu took 100' \
  'Rudro,Arnob,Rabbi,Arman,Nihal,Biswarup,Shohag,Rizu,Sarwar,Saimun,Tarek,Naimun,Anzum'

add_match 2026-09-02 2040 '' \
  'Shohag,Arman,Rabbi,Arnob,Rizu:1,Anzum,Sarwar,Saber'

add_match 2026-09-08 2060 '' \
  'Shohag,Arman,Sadi,Rabbi,Arnob,Anzum,Rizu,Amir,Tarek,Kawchar'

add_match 2026-09-12 4060 '' \
  'Arnob,Sadi,Arman,Imtiaz,Shohag,Sarwar,Anzum,Rizu,Kawchar,Maruf,Shuvo,Munjib'

add_match 2026-09-16 2610 'Bill 2110 + 500 moved from Sep 26' \
  'Shohag,Rizu,Rashid,Maruf,Arman,Biswarup,Rabbi,Nihal,Arnob,Sarwar,Saber,Sadman'

add_match 2026-09-26 3600 'Bill 4100 - 500 moved to Sep 16' \
  'Arnob,Biswarup,Rudro,Sadi,Shohag,Rizu,Sarwar,Kawchar,Alfeh'

add_match 2026-10-01 2020 '' \
  'Arnob,Rabbi,Sarwar,Alfeh,Anzum:1,Kawchar'

add_match 2026-10-03 4040 '' \
  'Arnob,Sadi,Nihal,Imtiaz,Amir,Shohag,Rizu,Sarwar,Anzum,Alfeh,Kawchar,Maruf,Rashid,Khaled,Porag'

add_match 2026-10-10 4000 '' \
  'Arnob,Rasel,Nihal,Imtiaz,Shohag,Anzum,Sarwar,Kawchar,Alfeh,Tarek,Prosanto,Rukonur,Saber'
