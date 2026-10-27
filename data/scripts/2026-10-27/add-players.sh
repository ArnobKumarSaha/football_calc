#!/usr/bin/env bash

set -euo pipefail

BASE_URL=${BASE_URL:-http://football.kubedb.cloud}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-changeme}


for name in Naimun Saimun Fahim Nihal Anzum Maruf Shuvo \
            Munjib Rashid Sadman Khaled Porag Prosanto Rukonur; do
  curl -sS --fail-with-body -X POST "$BASE_URL/api/players" \
    -H "Authorization: Bearer $ADMIN_PASSWORD" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"$name\"}"
  echo
done
