# ChatVault Kubernetes Deployment

## Day 3 Checklist

| Step | Task | Status |
|------|------|--------|
| 1 | K8s manifests (Deployment, Service, ConfigMap, Ingress) | ✅ Manifests ready |
| 2 | Deploy to kind | See below |
| 3 | Add Ingress | ✅ Ingress included |
| 4 | Test | Running in K8s |

---

## Quick Start (kind)

### 1. Create kind cluster with Ingress

```bash
# Create kind cluster
kind create cluster --config deployments/kind-config.yaml

# Install NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
kubectl wait --namespace ingress-nginx --for=condition=ready pod -l app.kubernetes.io/component=controller --timeout=90s
```

### 2. Build and load image

```bash
docker build -t chatvault:latest .
kind load docker-image chatvault:latest
```

### 3. Deploy

```bash
kubectl apply -k deployments/
```

### 4. Test

```bash
# Add host entry (or use 127.0.0.1)
echo "127.0.0.1 chatvault.local" | sudo tee -a /etc/hosts

# Health check
curl http://chatvault.local/health

# API
curl http://chatvault.local/api/v1/...
```

---

## Manifest Overview

| File | Purpose |
|------|---------|
| `namespace.yaml` | chatvault namespace |
| `configmap.yaml` | App port, DB name |
| `secret.yaml` | DB credentials, JWT key, DATABASE_URL |
| `postgres-pvc.yaml` | 5Gi persistent volume for DB |
| `postgres-deployment.yaml` | Postgres Deployment + Service |
| `app-deployment.yaml` | ChatVault Deployment + Service |
| `ingress.yaml` | Ingress (host: chatvault.local) |
| `kustomization.yaml` | Kustomize config |

---

## Alternative: Port-forward (no Ingress)

```bash
kubectl port-forward svc/chatvault-app 8080:80 -n chatvault
curl http://localhost:8080/health
```
