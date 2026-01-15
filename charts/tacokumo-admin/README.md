# Tacokumo Admin Helm Chart

This Helm chart deploys the Tacokumo Admin system, which consists of:
- **Admin API**: Go-based REST API with GitHub OAuth authentication
- **Admin UI**: Next.js-based frontend application
- **External Dependencies**: PostgreSQL and Redis/Valkey services

## Prerequisites

Before deploying this chart, ensure you have:

1. **External PostgreSQL Database**:
   - Database: `tacokumo_admin_db`
   - User: `admin_api` with appropriate permissions

2. **External Redis/Valkey Instance**:
   - Accessible from the Kubernetes cluster
   - Optional password authentication

3. **GitHub OAuth Application**:
   - Client ID and Client Secret
   - Callback URL configured: `https://your-domain.com/v1alpha1/auth/callback`

4. **Kubernetes Cluster**:
   - Ingress controller (nginx recommended)
   - Optional: TLS certificates for HTTPS

## Configuration

### Required Secrets

Create the following secrets before deploying:

```bash
# GitHub OAuth credentials
kubectl create secret generic tacokumo-admin-github-oauth \
  --from-literal=clientId="your-github-client-id" \
  --from-literal=clientSecret="your-github-client-secret"

# Database credentials
kubectl create secret generic tacokumo-admin-db-credentials \
  --from-literal=password="your-database-password"

# Optional: Redis credentials (if authentication required)
kubectl create secret generic tacokumo-admin-redis-credentials \
  --from-literal=password="your-redis-password"
```

### Values Configuration

Update `values.yaml` with your specific configuration:

```yaml
# External service endpoints
global:
  externalServices:
    postgresql:
      host: "your-postgresql-host"
      port: 5432
      database: "tacokumo_admin_db"
      username: "admin_api"
    redis:
      host: "your-redis-host"
      port: 6379

# GitHub OAuth settings
github:
  oauth:
    callbackUrl: "https://your-domain.com/v1alpha1/auth/callback"
    allowedOrgs: "your-org-1,your-org-2"

# Ingress configuration
ingress:
  hosts:
    - host: your-domain.com
      paths:
        - path: /
          pathType: Prefix
          service: tacokumo-admin-ui
          port: 3000
        - path: /api
          pathType: Prefix
          service: tacokumo-admin-api
          port: 8080
  tls:
    - secretName: your-tls-secret
      hosts:
        - your-domain.com
```

## Deployment

1. **Install or upgrade the chart**:
   ```bash
   helm upgrade --install tacokumo-admin ./charts/tacokumo-admin \
     --namespace tacokumo-admin \
     --create-namespace \
     --values your-values.yaml
   ```

2. **Verify deployment**:
   ```bash
   kubectl get pods -n tacokumo-admin
   kubectl get services -n tacokumo-admin
   kubectl get ingress -n tacokumo-admin
   ```

3. **Check health endpoints**:
   ```bash
   # API health
   curl https://your-domain.com/api/v1alpha1/health/liveness

   # UI health
   curl https://your-domain.com/api/health
   ```

## Service Communication

- **External Access**: `https://your-domain.com`
  - UI: `/` (served by Next.js on port 3000)
  - API: `/api/*` (proxied to Go service on port 8080)
- **Internal Communication**: UI → API via `http://tacokumo-admin-api:8080`
- **Authentication Flow**: API handles GitHub OAuth, UI displays login/logout

