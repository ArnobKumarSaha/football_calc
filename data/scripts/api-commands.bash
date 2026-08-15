BASE_URL=http://football.kubedb.cloud        # or http://localhost:8080
ADMIN_PASSWORD=changeme

# Check auth
curl -sS -o /dev/null -w '%{http_code}\n' "$BASE_URL/api/auth" \
  -H "Authorization: Bearer $ADMIN_PASSWORD"

# List players
curl -s "$BASE_URL/api/players" | jq -r '.[].name'

# Cleanup
curl -sS -X POST "$BASE_URL/api/admin/cleanup" \
  -H "Authorization: Bearer $ADMIN_PASSWORD"
