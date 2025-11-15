# Deployment Guide

## Table of Contents
1. [Local Development](#local-development)
2. [Docker Deployment](#docker-deployment)
3. [Kubernetes Deployment](#kubernetes-deployment)
4. [Production Checklist](#production-checklist)
5. [Monitoring](#monitoring)
6. [Troubleshooting](#troubleshooting)

---

## Local Development

### Prerequisites
- Go 1.22 or later
- Make (optional)

### Quick Start

1. **Download dependencies:**
```bash
go mod download
```

2. **Run the service:**
```bash
go run ./cmd/service
```

Or use Make:
```bash
make run
```

3. **Test the service:**
```bash
curl http://localhost:8080/health
```

### Configuration

Edit `config.yaml` or use environment variables:

```bash
# Override server address
export SERVER_ADDR=":9090"

# Override JWT secret
export JWT_SECRET="production-secret-key"

# Override log level
export LOG_LEVEL="debug"

# Run with custom config
go run ./cmd/service
```

---

## Docker Deployment

### Build Docker Image

Create `Dockerfile`:
```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service ./cmd/service

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/service .
COPY config.yaml .

EXPOSE 8080

CMD ["./service"]
```

Build and run:
```bash
# Build image
docker build -t myservice:latest .

# Run container
docker run -p 8080:8080 myservice:latest

# Run with environment variables
docker run -p 8080:8080 \
  -e JWT_SECRET="production-secret" \
  -e LOG_LEVEL="info" \
  myservice:latest
```

### Docker Compose

Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  service:
    build: .
    ports:
      - "8080:8080"
    environment:
      - JWT_SECRET=${JWT_SECRET}
      - LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

Run:
```bash
docker-compose up -d
```

---

## Kubernetes Deployment

### ConfigMap

Create `k8s/configmap.yaml`:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: myservice-config
data:
  config.yaml: |
    server:
      addr: ":8080"
      read_timeout: 10s
      write_timeout: 10s
      shutdown_timeout: 30s
    logging:
      level: info
      format: json
    cors:
      allowed_origins:
        - "*"
      allowed_methods:
        - GET
        - POST
        - PUT
        - DELETE
      allowed_headers:
        - Content-Type
        - Authorization
    rate_limit:
      requests_per_second: 100
      burst: 10
    jwt:
      secret: ${JWT_SECRET}
      expiration: 24h
```

### Secret

Create `k8s/secret.yaml`:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: myservice-secrets
type: Opaque
stringData:
  jwt-secret: "your-production-secret-key-change-this"
```

Create secret:
```bash
kubectl create secret generic myservice-secrets \
  --from-literal=jwt-secret="$(openssl rand -base64 32)"
```

### Deployment

Create `k8s/deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myservice
  labels:
    app: myservice
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myservice
  template:
    metadata:
      labels:
        app: myservice
    spec:
      containers:
      - name: myservice
        image: myservice:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: myservice-secrets
              key: jwt-secret
        - name: LOG_LEVEL
          value: "info"
        volumeMounts:
        - name: config
          mountPath: /root/config.yaml
          subPath: config.yaml
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 3
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
      volumes:
      - name: config
        configMap:
          name: myservice-config
```

### Service

Create `k8s/service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: myservice
  labels:
    app: myservice
spec:
  type: ClusterIP
  selector:
    app: myservice
  ports:
  - port: 8080
    targetPort: 8080
    protocol: TCP
    name: http
```

### Ingress

Create `k8s/ingress.yaml`:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myservice
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - host: myservice.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: myservice
            port:
              number: 8080
```

### Deploy to Kubernetes

```bash
# Apply all resources
kubectl apply -f k8s/

# Check deployment status
kubectl get deployments
kubectl get pods
kubectl get services

# View logs
kubectl logs -l app=myservice --tail=100 -f

# Scale deployment
kubectl scale deployment myservice --replicas=5

# Rolling update
kubectl set image deployment/myservice myservice=myservice:v2

# Rollback
kubectl rollout undo deployment/myservice
```

---

## Production Checklist

### Security
- [ ] Change default JWT secret
- [ ] Use HTTPS/TLS certificates
- [ ] Enable rate limiting
- [ ] Implement authentication on all protected endpoints
- [ ] Validate all user inputs
- [ ] Use secrets management (Vault, AWS Secrets Manager)
- [ ] Enable CORS only for trusted origins
- [ ] Implement request size limits
- [ ] Add security headers (X-Frame-Options, CSP, etc.)

### Performance
- [ ] Configure appropriate resource limits
- [ ] Enable HTTP/2
- [ ] Add caching layer (Redis)
- [ ] Optimize database queries
- [ ] Use connection pooling
- [ ] Enable gzip compression
- [ ] Configure CDN for static assets

### Observability
- [ ] Set up centralized logging (ELK, Loki)
- [ ] Configure metrics collection (Prometheus)
- [ ] Set up alerting (Alertmanager)
- [ ] Add distributed tracing (Jaeger, Zipkin)
- [ ] Create dashboards (Grafana)
- [ ] Implement health checks
- [ ] Add performance monitoring (APM)

### Reliability
- [ ] Configure auto-scaling (HPA)
- [ ] Set up load balancing
- [ ] Implement circuit breakers
- [ ] Add retry logic with exponential backoff
- [ ] Configure graceful shutdown
- [ ] Set up database backups
- [ ] Implement disaster recovery plan
- [ ] Test failover scenarios

### Configuration
- [ ] Use environment-specific configs
- [ ] Externalize all secrets
- [ ] Document all configuration options
- [ ] Validate configuration at startup
- [ ] Version control configuration files

---

## Monitoring

### Prometheus Metrics

The service exposes metrics at `/metrics`:

```promql
# Request rate
rate(http_requests_total[5m])

# Error rate
rate(http_requests_total{status=~"5.."}[5m])

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Active requests
http_active_requests
```

### Grafana Dashboard

Import Grafana dashboard JSON:
```json
{
  "dashboard": {
    "title": "Microservice Metrics",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])"
          }
        ]
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])"
          }
        ]
      }
    ]
  }
}
```

### Alerting Rules

Create `prometheus-rules.yaml`:
```yaml
groups:
- name: microservice
  interval: 30s
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value }} req/s"

  - alert: HighLatency
    expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High latency detected"
      description: "P95 latency is {{ $value }}s"
```

---

## Troubleshooting

### Service Won't Start

**Check logs:**
```bash
# Docker
docker logs <container-id>

# Kubernetes
kubectl logs -l app=myservice
```

**Common issues:**
- Missing JWT secret → Set `JWT_SECRET` env var
- Port already in use → Change `SERVER_ADDR`
- Invalid config → Check `config.yaml` syntax

### High Memory Usage

**Check metrics:**
```bash
curl http://localhost:8080/metrics | grep go_memstats
```

**Solutions:**
- Reduce connection pool size
- Add memory limits in Kubernetes
- Enable garbage collection tuning

### High Latency

**Check request duration:**
```bash
curl http://localhost:8080/metrics | grep http_request_duration
```

**Solutions:**
- Add caching
- Optimize database queries
- Scale horizontally
- Add CDN

### Rate Limiting Issues

**Check configuration:**
```yaml
rate_limit:
  requests_per_second: 100  # Increase if needed
  burst: 10                  # Increase burst capacity
```

### Database Connection Errors

**Check readiness probe:**
```bash
curl http://localhost:8080/ready
```

**Solutions:**
- Verify database is running
- Check connection string
- Increase connection timeout
- Add connection retry logic

---

## Useful Commands

```bash
# Build
make build

# Run locally
make run

# Run tests
make test

# Run with race detector
make test-race

# Format code
make fmt

# Clean build artifacts
make clean

# View logs (Docker)
docker logs -f <container-id>

# View logs (Kubernetes)
kubectl logs -f deployment/myservice

# Port forward (Kubernetes)
kubectl port-forward service/myservice 8080:8080

# Exec into pod
kubectl exec -it <pod-name> -- /bin/sh

# View resource usage
kubectl top pods -l app=myservice
```

---

## Support

For issues and questions:
- Check logs first
- Review configuration
- Test with curl commands (see TEST.md)
- Check metrics endpoint
- Review health/ready endpoints
