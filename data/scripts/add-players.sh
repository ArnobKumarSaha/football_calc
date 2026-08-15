#!/usr/bin/env bash

set -euo pipefail

BASE_URL=${BASE_URL:-http://football.kubedb.cloud}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-changeme}


for name in Imtiaz Mamun Lotifur Mennal Avishek Saber Nirjhor \
            Kawchar Sabbir Farzine Meraj Sarwar Shakil Tarek Alfeh; do
  curl -sS -X POST "$BASE_URL/api/players" \
    -H "Authorization: Bearer $ADMIN_PASSWORD" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"$name\"}"
  echo
done
