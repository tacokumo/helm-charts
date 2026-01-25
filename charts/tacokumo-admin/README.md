# Tacokumo Admin Helm Chart

This Helm chart deploys the Tacokumo Admin system, a unified Go application built with Echo framework and Go templates for the UI.

## Architecture

The admin application is a single monolithic service that includes:
- **Web UI**: Server-rendered pages using Go templates
- **REST API**: Backend API endpoints
- **GitHub OAuth**: Authentication via GitHub organization membership
- **External Dependencies**: PostgreSQL and Redis/Valkey services

## Prerequisites

Before deploying this chart, ensure you have:

1. **External PostgreSQL Database**:
   - Database: `tacokumo_admin_db`
   - User with appropriate permissions

2. **External Redis/Valkey Instance**:
   - Accessible from the Kubernetes cluster
   - Optional password authentication

3. **GitHub OAuth Application**:
   - Client ID and Client Secret
   - Callback URL configured: `https://your-domain.com/auth/callback`

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
      username: "admin"
      initialConnRetry: 10
    redis:
      host: "your-redis-host"
      port: 6379
      db: 0
      initialConnRetry: 10

# GitHub OAuth settings
github:
  oauth:
    callbackUrl: "https://your-domain.com/auth/callback"
    org: "your-github-org"
  session:
    ttl: "24h"
    cookieSecure: true
  frontendUrl: "https://your-domain.com"

# Admin application settings
admin:
  baseDomain: "tacokumo.dev"
  config:
    logLevel: "info"
    opentelemetry:
      enabled: false
      serviceName: "tacokumo-admin"
      tracesExporter: "otlp"
      otlpEndpoint: "http://otel-collector:4317"
      otlpProtocol: "grpc"

# Ingress configuration
ingress:
  enabled: true
  hosts:
    - host: your-domain.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: your-tls-secret
      hosts:
        - your-domain.com
```

### Key Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `admin.replicaCount` | Number of replicas | `1` |
| `admin.baseDomain` | Base domain for projects | `"tacokumo.dev"` |
| `admin.config.logLevel` | Log level (debug, info, warn, error) | `"info"` |
| `admin.config.opentelemetry.enabled` | Enable OpenTelemetry tracing | `false` |
| `admin.livenessProbe.path` | Liveness probe endpoint | `/healthz` |
| `admin.readinessProbe.path` | Readiness probe endpoint | `/readyz` |
| `github.oauth.org` | GitHub organization for access control | `""` |
| `github.session.ttl` | Session TTL duration | `"24h"` |
| `github.session.cookieSecure` | Enable secure cookies | `true` |

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
   # Liveness check
   curl https://your-domain.com/healthz

   # Readiness check
   curl https://your-domain.com/readyz
   ```

## Service Communication

- **External Access**: `https://your-domain.com`
  - All routes served by the unified admin service on port 8080
- **Authentication Flow**: GitHub OAuth handled internally
- **Session Management**: Redis-backed sessions
