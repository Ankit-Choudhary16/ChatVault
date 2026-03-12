Fill all the neccessary infomation in the .env file 

### Run Locally

```bash
# Terminal 1: Start PostgreSQL
make db-up
export DATABASE_URL="postgres://postgres:password@localhost:5432/chatvault?sslmode=disable"

# Terminal 2: Run the app
cp .env.example .env
make run

# Test health endpoint
curl http://localhost:8080/health
```

---

### To run it on docker 

```bash
# build the image 
docker build -t chatvault:latest .

#run it 
docker compose up 






### 1. Create kind cluster with Ingress

```bash
# Create kind cluster
kind create cluster --config deployments/kind-config.yaml
or 
make create-cluster

# Install NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
or 
make ingress-controller
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
