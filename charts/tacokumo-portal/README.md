# tacokumo-portal

A Helm chart for the Portal-API application.

## Description

This chart deploys the Portal-API application, which provides API endpoints for managing portal resources.

## Requirements

- Kubernetes 1.19+
- Helm 3.0+

## Configuration

The following table lists the configurable parameters of the tacokumo-portal chart.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `api.portalName` | Portal namespace name (REQUIRED) | `"default-portal"` |
| `api.logLevel` | Logging level (debug, info, warn, error) | `"info"` |
| `api.image.repository` | Container image repository | `ghcr.io/tacokumo/portal-api` |
| `api.image.tag` | Container image tag | `"latest"` |
| `api.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `api.hpa.enabled` | Enable HPA | `true` |
| `api.hpa.minReplicas` | Minimum replicas | `1` |
| `api.hpa.maxReplicas` | Maximum replicas | `3` |
| `api.hpa.targetMemoryUtilizationPercentage` | Target memory utilization | `80` |
| `api.service.enabled` | Enable Service | `true` |
| `api.service.type` | Service type | `ClusterIP` |
| `api.service.port` | Service port | `1323` |
| `api.rbac.create` | Create RBAC resources | `true` |
| `api.serviceAccount.create` | Create ServiceAccount | `true` |
| `api.serviceAccount.name` | ServiceAccount name | `"portal-api"` |

## RBAC Permissions

The chart creates a namespace-scoped Role with the following permissions:

- `tacokumo.github.io/applications`: get, list, create, watch
- `core/secrets`: get, create, update

## Health Endpoints

- Liveness: `/health/liveness`
- Readiness: `/health/readiness`
