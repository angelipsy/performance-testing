# FastAPI vs Go Service Performance Testing

Performance comparison project for testing FastAPI (Python) and Go services under various load conditions using Kubernetes and k6.

## Overview

This project benchmarks two microservices implementations:
- **FastAPI Service**: Python-based REST API built with FastAPI
- **Go Service**: Go-based REST API

The services are deployed on a local Kubernetes cluster (kind) and tested using k6 load testing jobs. The comparison covers multiple scenarios including health checks, CPU-intensive operations, I/O operations, and JSON processing.

## Project Structure

```
.
├── apps/
│   ├── fastapi/          
│   └── golang/           
├── k8s/                  
│   ├── fastapi-service/  
│   ├── go-service/       
│   ├── prometheus/       
│   ├── k6-job.yaml       
│   └── k6-job-stress.yaml 
├── test-load-results/    
├── test-stress-results/  
├── scripts/
|   ├── *.sh              
└── README.md
```

## Prerequisites

- Ubuntu/Debian-based Linux system
- Sudo privileges for installing Docker and Kubernetes tools

## Setup and Installation

### 1. Install Dependencies

On first run, install all required dependencies:

```sh
./scripts/installations.sh
```

This script installs:
- Docker Engine with Compose plugin
- kubectl (Kubernetes CLI)
- kind (Kubernetes in Docker)

**Note**: You may need to re-login or run `newgrp docker` after installation to use Docker without sudo.

### 2. Create Kubernetes Cluster

Create a local kind cluster:

```sh
./scripts/cluster-creation.sh
```

### 3. Build and Load Docker Images

Build both service Docker images and load them into the kind cluster:

```sh
./scripts/load-docker-images-in-cluster.sh
```

### 4. Deploy Applications

Deploy both services to the Kubernetes cluster:

```sh
./scripts/deploy-apps.sh
```

Wait for all pods to be created and running. You can check the status with:

```sh
kubectl -n poc-compare get pods
```

## Running Tests

### Batch Testing

Run all load tests in sequence:

```sh
./scripts/run-tests.sh
```

This executes k6 load tests for both services and saves results to `test-load-results/`.

Run all stress tests in sequence:

```sh
./scripts/run-stress-tests.sh
```

This executes k6 stress tests for both services and saves results to `test-stress-results/`.

### Individual Test Execution

Run a single test and view results:

```sh
kubectl -n poc-compare logs job/k6-golang-base -f > test-load-results/k6-golang-health-results.txt
```

Available test jobs:
- `k6-golang-base` / `k6-fastapi-base` - Health endpoint tests
- `k6-golang` / `k6-fastapi` - CPU-intensive tests
- `k6-golang-io` / `k6-fastapi-io` - I/O operation tests
- `k6-golang-json` / `k6-fastapi-json` - JSON processing tests

For stress tests, use the `k6-stress-*` prefix:
- `k6-stress-golang-base` / `k6-stress-fastapi-base`
- `k6-stress-golang` / `k6-stress-fastapi`
- `k6-stress-golang-io` / `k6-stress-fastapi-io`
- `k6-stress-golang-json` / `k6-stress-fastapi-json`

## Test Results

Test results are saved in:
- `test-load-results/` - Load test outputs
- `test-stress-results/` - Stress test outputs

Each file contains k6 metrics including:
- Request rates (requests/sec)
- Response times (min, avg, max, p95, p99)
- Success/failure rates
- Virtual users (VUs)

## Monitoring

The project includes Prometheus and Grafana for monitoring:

### Port Forwarding

Forward FastAPI service:
```sh
./scripts/port-forward-fastapi.sh
```

Forward Go service:
```sh
./scripts/port-forward-golang.sh
```

### Grafana Dashboard

A pre-configured Grafana dashboard is available in `grafana-dashboard.json` for visualizing metrics.

## Local Development

Run both services locally with Docker Compose:

```sh
docker compose up --build
```

Test endpoints:
```sh
curl localhost:8001/health    # FastAPI health check
curl localhost:8002/health    # Go health check
curl "localhost:8001/cpu?n=30"  # FastAPI CPU test
curl "localhost:8002/cpu?n=30"  # Go CPU test
```

## Additional Scripts

- `apply-hpa-change.sh` - Apply Horizontal Pod Autoscaler changes
- `port-forward-fastapi.sh` - Port forward to FastAPI service
- `port-forward-golang.sh` - Port forward to Go service

## Configuration Files

- `docker-compose.yml` - Local development setup
- `prometheus.yml` - Prometheus configuration
- `grafana-dashboard.json` - Grafana dashboard definition

## Cleanup

To delete the kind cluster:

```sh
kind delete cluster
```

## Notes

- The namespace used for all deployments is `poc-compare`
- Make sure to wait for pods to reach `Running` state before executing tests
- Test results are appended to files, so clean up old results if needed
