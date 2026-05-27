# Kubernetes Setup

## Prerequisites

- K8s cluster with Gateway API CRDs installed
- A `Gateway` named `lb-gateway` in namespace `gateway`
- External managed PostgreSQL instance

## 1. Build and push image

```bash
make push
# pushes ghcr.io/arnobkumarsaha/football-calc:latest for linux/amd64,linux/arm64
```

To use a custom tag:
```bash
make push TAG=v1.0.0
```

## 2. Configure the secret

Edit `deploy/secret.yaml` with real values before applying:

```yaml
stringData:
  ADMIN_PASSWORD: "<your-password>"
  DATABASE_URL: "postgres://user:password@host:5432/football_calc?sslmode=require"
```

> Do not commit real credentials. Use `kubectl create secret` or a secrets manager in production.

```bash
kubectl create secret generic football-calc-secret \
  --namespace football-calc \
  --from-literal=ADMIN_PASSWORD=<password> \
  --from-literal=DATABASE_URL="postgres://user:pass@host:5432/football_calc?sslmode=require"
```

## 3. Apply manifests

```bash
kubectl apply -f deploy/namespace.yaml
kubectl apply -f deploy/secret.yaml      # or use kubectl create secret above
kubectl apply -f deploy/deployment.yaml
kubectl apply -f deploy/service.yaml
kubectl apply -f deploy/httproute.yaml
```

Or all at once (if using kustomize):
```bash
kubectl apply -k deploy/
```

## 4. Verify

```bash
kubectl -n football-calc get pods
kubectl -n football-calc logs deploy/football-calc

# Health checks via LB
curl http://<lb-ip>/healthz
curl http://<lb-ip>/readyz
```

## Manifests Summary

| File | Resource | Notes |
|------|----------|-------|
| `namespace.yaml` | Namespace `football-calc` | |
| `secret.yaml` | Secret `football-calc-secret` | `ADMIN_PASSWORD`, `DATABASE_URL` |
| `deployment.yaml` | Deployment 1 replica | envFrom secret, liveness+readiness probes |
| `service.yaml` | Service ClusterIP :8080 | |
| `httproute.yaml` | HTTPRoute | routes `/` → service:8080 via `lb-gateway` |

## Resource Limits

```
requests: cpu=50m, memory=64Mi
limits:   cpu=200m, memory=128Mi
```

## Update image

```bash
kubectl -n football-calc set image deployment/football-calc football-calc=ghcr.io/arnobkumarsaha/football-calc:<new-tag>
kubectl -n football-calc rollout status deployment/football-calc
```
