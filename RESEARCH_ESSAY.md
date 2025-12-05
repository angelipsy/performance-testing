# Performance Comparison of FastAPI and Golang Services: A Comprehensive Study

## Abstract

The selection of an appropriate backend technology stack is a critical decision that significantly impacts the scalability, reliability, and operational costs of modern web applications. This research project aims to conduct a systematic performance comparison between FastAPI (Python) and Golang to identify the optimal platform for developing high-performance backend services.

FastAPI was selected as the Python representative due to its position as one of the fastest Python web frameworks available. Built on Starlette and leveraging asynchronous capabilities through Python's asyncio, FastAPI has gained widespread adoption for its combination of high performance, automatic API documentation, type safety through Pydantic, and developer-friendly syntax. The framework represents the pinnacle of what Python can achieve in web service performance while maintaining the language's renowned ease of development.

Golang, conversely, was designed from the ground up at Google to address the challenges of building large-scale, concurrent services. Its compiled nature, built-in concurrency primitives (goroutines and channels), and minimal runtime overhead make it a natural candidate for high-throughput, low-latency applications.

The motivation for this comparative study stems from the need to make data-driven decisions when architecting backend services. While anecdotal evidence and benchmarks exist, controlled performance testing under realistic workloads with identical resource constraints provides the empirical foundation necessary for informed technology selection. This project establishes a reproducible testing methodology to evaluate both platforms across multiple dimensions: baseline throughput, CPU-intensive operations, I/O-bound tasks, and JSON serialization performance.

## 1. Introduction: Key Differences Between FastAPI and Golang

### 1.1 Language Paradigms and Execution Models

**FastAPI (Python)** is an interpreted, dynamically-typed language framework that leverages modern Python features:
- **Asynchronous Execution**: Built on Python's asyncio event loop, enabling non-blocking I/O operations
- **ASGI-based**: Uses the Asynchronous Server Gateway Interface through Uvicorn
- **Dynamic Typing with Runtime Validation**: Pydantic models provide runtime type checking and validation
- **GIL Limitations**: Python's Global Interpreter Lock restricts true CPU parallelism to a single thread per process

**Golang** is a statically-typed, compiled language with native concurrency support:
- **Compiled to Native Code**: Produces standalone binaries with no runtime dependencies
- **Goroutines**: Lightweight threads managed by the Go runtime, enabling massive concurrency with minimal overhead
- **Static Typing**: Compile-time type checking catches errors before deployment
- **Efficient Memory Management**: Garbage collection optimized for low-latency scenarios
- **True Parallelism**: Goroutines can execute across multiple CPU cores without GIL constraints

### 1.2 Concurrency Models

The fundamental difference lies in how each platform handles concurrent requests:

**FastAPI**: Uses an event loop with cooperative multitasking. Async/await syntax allows yielding control during I/O operations, but CPU-bound work blocks the event loop unless offloaded to thread/process pools.

**Golang**: Each HTTP request typically runs in its own goroutine. The Go runtime schedules thousands of goroutines across available CPU cores, providing both efficient I/O multiplexing and true parallel CPU execution.

### 1.3 Development and Deployment

**FastAPI** offers:
- Rapid development with automatic API documentation (OpenAPI/Swagger)
- Rich ecosystem of Python libraries
- Easier debugging and more verbose error messages
- Requires Python runtime and dependencies in deployment

**Golang** provides:
- Strong compile-time guarantees
- Single binary deployment with no external dependencies
- More verbose code with explicit error handling
- Faster compilation and startup times

## 2. Technologies

### 2.1 Core Application Frameworks

#### FastAPI (v0.115.0)
**What it is**: A modern, high-performance Python web framework for building APIs with automatic interactive documentation.

**What it's used for**: Building RESTful APIs and microservices with automatic request validation, serialization, and API documentation generation.

**Role in this project**: Serves as the Python-based service implementation, exposing four test endpoints (/health, /cpu, /io, /json) to evaluate different performance characteristics under load.

**Key dependencies**:
- **Uvicorn (v0.30.6)**: ASGI server that runs the FastAPI application
- **Gunicorn (v22.0.0)**: Production WSGI server for process management (optional for scaling)

#### Golang (v1.23.0)
**What it is**: An open-source programming language developed by Google, designed for building efficient, reliable, and scalable software.

**What it's used for**: Building high-performance services, microservices, and concurrent applications.

**Role in this project**: Implements an equivalent HTTP service using Go's standard library `net/http` package, exposing identical endpoints to FastAPI for direct performance comparison.

### 2.2 Monitoring and Observability

#### Prometheus
**What it is**: An open-source monitoring and alerting toolkit with a dimensional data model and powerful query language (PromQL).

**What it's used for**: Collecting time-series metrics from instrumented applications, storing them, and providing a query interface for analysis.

**Role in this project**:
- Scrapes metrics from both FastAPI and Golang services every 5 seconds
- Collects application-level metrics (request counts, duration histograms, HTTP status codes)
- Enables real-time monitoring during load and stress tests
- Provides data source for Grafana visualizations

**Configuration**: Located in `/k8s/prometheus/prometheus-configmap.yaml` with two scrape jobs:
- `fastapi` job targeting the FastAPI service
- `golang` job targeting the Golang service

#### Prometheus Client Libraries
**FastAPI integration**: Uses `prometheus-fastapi-instrumentator` (v7.0.0) which automatically instruments all endpoints with:
- Request duration histograms
- Request counts by method, path, and status code
- In-progress request gauges

**Golang integration**: Uses official `prometheus/client_golang` (v1.23.2) with custom middleware implementing:
- `http_requests_total`: Counter tracking requests by method, endpoint, and status
- `http_request_duration_seconds`: Histogram measuring request latency

#### Grafana
**What it is**: An open-source analytics and visualization platform that queries, visualizes, and alerts on metrics.

**What it's used for**: Creating dashboards to visualize time-series data from multiple sources.

**Role in this project**:
- Connects to Prometheus as a data source
- Provides real-time visualization of service performance during tests
- Displays metrics such as request rate, latency percentiles (p50, p90, p95, p99), error rates, and resource utilization
- Dashboard configuration stored in `grafana-dashboard.json`

### 2.3 Performance Testing

#### k6 (v0.47.0)
**What it is**: An open-source load testing tool built for developer experience, designed to test the performance of APIs, microservices, and websites.

**What it's used for**: Simulating realistic user traffic patterns with configurable virtual users (VUs) and generating detailed performance reports.

**Role in this project**:
- Executes load tests with ramping VU profiles
- Runs stress tests to identify breaking points
- Tests individual endpoints (/health, /cpu, /io, /json) in isolation
- Provides detailed metrics on request duration, failure rates, and throughput
- Configured via JavaScript test scripts with defined stages and thresholds

**Test Scripts**:
- `load.js`: Progressive ramp from 20 to 60 VUs over 2 minutes
- `stress.js`: Aggressive ramp from 10 to 30 VUs over 2.5 minutes (adjusted for single-pod testing)

### 2.4 Container Orchestration

#### Kubernetes
**What it is**: An open-source container orchestration platform for automating deployment, scaling, and management of containerized applications.

**What it's used for**: Managing containerized workloads and services with declarative configuration.

**Role in this project**:
- Orchestrates both FastAPI and Golang service deployments
- Manages resource allocation (CPU/memory limits and requests)
- Runs k6 load testing jobs in the same cluster for network isolation
- Hosts Prometheus and Grafana monitoring stack
- Provides service discovery and load balancing

**Namespace**: All resources deployed in the `poc-compare` namespace for isolation

#### Docker
**What it is**: A platform for developing, shipping, and running applications in containers.

**What it's used for**: Packaging applications with their dependencies into standardized units.

**Role in this project**:
- Builds container images for FastAPI service (Dockerfile in `apps/fastapi/`)
- Builds container images for Golang service (Dockerfile in `apps/golang/`)
- Enables consistent deployment across environments
- Images tagged as `poc-fastapi:latest` and `poc-golang:latest`

### 2.5 Summary of Technology Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Python Framework | FastAPI | 0.115.0 | API service implementation |
| ASGI Server | Uvicorn | 0.30.6 | FastAPI application server |
| Language | Golang | 1.23.0 | Alternative service implementation |
| Load Testing | k6 | 0.47.0 | Performance testing and benchmarking |
| Metrics Collection | Prometheus | Latest | Time-series metrics database |
| Visualization | Grafana | Latest | Metrics dashboards and analysis |
| Orchestration | Kubernetes | Latest | Container management and scaling |
| Containerization | Docker | Latest | Application packaging |

## 3. Performance Testing Solution: Design and Implementation

### 3.1 Testing Architecture

The performance testing solution is implemented entirely within a Kubernetes cluster to ensure:
- Network isolation and consistent network latency
- Reproducible test conditions
- Realistic deployment environment
- Easy resource constraint enforcement

**Architecture components**:
1. **Target Services**: FastAPI and Golang deployments running as Kubernetes Deployments with associated Services
2. **Load Generators**: k6 runs as Kubernetes Jobs, executing test scripts against in-cluster service endpoints
3. **Monitoring Stack**: Prometheus and Grafana deployed in the same namespace, continuously scraping metrics from both services
4. **Resource Isolation**: Each service runs with identical CPU and memory limits

### 3.2 Test Endpoints

Four endpoints were implemented identically in both services to test different performance characteristics:

#### 3.2.1 `/health` - Baseline Throughput
**Purpose**: Establish baseline HTTP handling overhead with minimal application logic.

**Implementation**:
- Returns plain text "ok" response
- No computation, no I/O, no serialization
- Measures pure HTTP request/response handling efficiency

**What it tells us**: The raw throughput capacity of each framework when overhead is minimized. Differences here highlight the performance of the HTTP stack, routing, and framework machinery.

#### 3.2.2 `/cpu` - CPU-Intensive Operations
**Purpose**: Test performance under CPU-bound workloads.

**Implementation**:
- Performs 100 iterations of SHA-256 hashing on benchmark data
- FastAPI: Uses hashlib.sha256() in synchronous code
- Golang: Uses crypto/sha256 package

**What it tells us**: How each platform handles CPU-intensive work:
- FastAPI's GIL limitations: CPU work blocks the event loop, potentially impacting concurrent request handling
- Golang's ability to parallelize CPU work across multiple goroutines on separate cores
- The efficiency of each language's cryptographic implementations

#### 3.2.3 `/io` - I/O-Bound Operations
**Purpose**: Evaluate performance with I/O wait times (simulating database queries, external API calls).

**Implementation**:
- Sleeps for 50 milliseconds to simulate I/O latency
- FastAPI: Uses `asyncio.sleep()` (non-blocking)
- Golang: Uses `time.Sleep()` in a goroutine

**What it tells us**: How efficiently each platform handles concurrent I/O operations:
- FastAPI's strength: Async/await allows the event loop to handle other requests during sleep
- Golang's concurrency: Each request in its own goroutine naturally handles I/O waiting
- The overhead of context switching and concurrency management

#### 3.2.4 `/json` - Serialization Performance
**Purpose**: Test JSON serialization/deserialization performance with realistic payloads.

**Implementation**:
- Creates an object with nested data: name, value, and an array of 5000 integers
- Returns the object as JSON
- Tests framework serialization efficiency

**What it tells us**:
- The efficiency of each platform's JSON encoding libraries
- Memory allocation patterns during serialization
- Framework overhead in response generation

### 3.3 Test Types and Their Purpose

#### 3.3.1 Load Testing
**Definition**: Evaluating system behavior under expected and peak normal conditions.

**Configuration** (k6-job.yaml - load.js):
```javascript
stages: [
  { duration: '30s', target: 20 },   // Ramp-up: 0 → 20 VUs
  { duration: '60s', target: 60 },   // Peak load: 20 → 60 VUs sustained
  { duration: '30s', target: 0  },   // Ramp-down: 60 → 0 VUs
]
```

**Total duration**: 2 minutes
**Peak VUs**: 60 concurrent virtual users
**Request interval**: 100ms sleep between requests (0.1s)

**Thresholds**:
- `http_req_failed: ['rate<0.01']` - Less than 1% request failures acceptable
- `http_req_duration: ['p(90)<2000']` - 90th percentile response time under 2 seconds

**What it tells us**:
- **Steady-state performance**: How each service performs under sustained expected load
- **Latency characteristics**: Response time distribution (p50, p90, p95, p99)
- **Throughput capacity**: Requests per second each service can handle
- **Stability**: Whether performance degrades over time under constant load
- **Resource efficiency**: CPU and memory consumption at various load levels

**Use case**: Validates that the service can handle normal production traffic volumes with acceptable latency and no errors.

#### 3.3.2 Stress Testing
**Definition**: Pushing the system beyond normal capacity to identify breaking points and failure modes.

**Configuration** (k6-job-stress.yaml - stress.js):
```javascript
stages: [
  { duration: '30s', target: 10 },   // Initial ramp: 0 → 10 VUs
  { duration: '30s', target: 20 },   // Increase: 10 → 20 VUs
  { duration: '1m',  target: 30 },   // Peak stress: 20 → 30 VUs sustained
  { duration: '30s', target: 0 },    // Ramp-down: 30 → 0 VUs
]
```

**Total duration**: 2.5 minutes
**Peak VUs**: 30 concurrent virtual users (adjusted for single-replica constraint)
**Request interval**: 100ms sleep between requests

**Thresholds**:
- `http_req_failed: ['rate<0.05']` - Less than 5% request failures acceptable (more lenient than load tests)
- `http_req_duration: ['p(90)<5000']` - 90th percentile response time under 5 seconds (higher tolerance)

**What it tells us**:
- **Breaking point**: At what load does the service start failing?
- **Failure modes**: How does the service fail? (timeouts, errors, crashes)
- **Degradation patterns**: Does latency increase linearly or exponentially under stress?
- **Recovery behavior**: Can the service recover when load decreases?
- **Resource saturation**: Which resource (CPU, memory, goroutines/event loop) becomes the bottleneck?

**Use case**: Identifies the maximum capacity and helps establish capacity planning guidelines. Also reveals how the service behaves during traffic spikes or DDoS scenarios.

### 3.4 Resources Assigned to Services

Both services are deployed with identical resource constraints to ensure fair comparison:

```yaml
resources:
  requests:
    cpu: "500m"      # 0.5 CPU cores guaranteed
    memory: "256Mi"  # 256 MiB guaranteed
  limits:
    cpu: "500m"      # Maximum 0.5 CPU cores
    memory: "256Mi"  # Maximum 256 MiB
```

**Rationale for these constraints**:

1. **CPU Limit (500m / 0.5 cores)**:
   - **Realistic constraint**: Many production services run on fractional CPU allocations in containerized environments
   - **Cost-effective**: Represents a typical small-to-medium service allocation in cloud environments
   - **Highlights differences**: Limited CPU amplifies the performance differences between interpreted and compiled languages
   - **Concurrency testing**: With limited CPU, we can observe how each platform manages concurrent request scheduling
   - **Single worker**: FastAPI configured with `UVICORN_WORKERS=1` to match Golang's single-process model

2. **Memory Limit (256Mi)**:
   - **Sufficient for testing**: Adequate for both services to handle test workloads without OOM issues
   - **Realistic**: Represents a lean microservice allocation
   - **Prevents memory from being the bottleneck**: Ensures tests measure CPU and concurrency performance, not memory constraints
   - **Monitoring focus**: Allows observation of memory allocation patterns and GC behavior

3. **Requests = Limits (guaranteed resources)**:
   - **Consistent performance**: Ensures Kubernetes doesn't throttle the pods
   - **Fair comparison**: Both services get exactly the same guaranteed resources
   - **Eliminates variability**: No resource contention with other pods

4. **Single Replica**:
   - **Controlled testing**: Both deployments run with `replicas: 1` to test single-instance performance
   - **HPA disabled during testing**: HorizontalPodAutoscaler configured but minimum replicas set to 1
   - **Fair comparison**: Tests the performance of a single instance of each platform, not their scaling characteristics

### 3.5 Test Configuration and Virtual Users

#### Why These VU Counts?

**Load Test (60 VUs max)**:
- **Single pod, limited CPU**: With 0.5 CPU cores, 60 concurrent users provides sufficient concurrency to saturate the CPU
- **0.1s sleep**: Each VU makes ~10 requests per second, so 60 VUs ≈ 600 requests/second potential
- **Realistic load**: Represents a moderate production load for a microservice
- **Observes concurrency handling**: Enough to see how each platform manages concurrent requests under CPU constraints

**Stress Test (30 VUs max - adjusted)**:
- **Originally planned higher**: The stress test was designed to push harder, but adjusted for single-pod testing
- **Focus on CPU saturation**: 30 VUs with CPU-intensive endpoints push the 0.5 core limit
- **Sustained stress**: 1-minute hold at peak to observe degradation over time
- **Failure identification**: Aggressive enough to potentially trigger failures or significant latency degradation

#### Test Matrix

The project runs a comprehensive matrix of tests:

| Test Type | Service | Endpoint | Purpose |
|-----------|---------|----------|---------|
| Load | FastAPI | /health | Baseline HTTP performance |
| Load | Golang | /health | Baseline HTTP performance |
| Load | FastAPI | /cpu | CPU-bound performance |
| Load | Golang | /cpu | CPU-bound performance |
| Load | FastAPI | /io | I/O concurrency handling |
| Load | Golang | /io | I/O concurrency handling |
| Load | FastAPI | /json | JSON serialization performance |
| Load | Golang | /json | JSON serialization performance |
| Stress | FastAPI | /health | Baseline stress tolerance |
| Stress | Golang | /health | Baseline stress tolerance |
| Stress | FastAPI | /cpu | CPU stress breaking point |
| Stress | Golang | /cpu | CPU stress breaking point |
| Stress | FastAPI | /io | I/O stress handling |
| Stress | Golang | /io | I/O stress handling |
| Stress | FastAPI | /json | Serialization under stress |
| Stress | Golang | /json | Serialization under stress |

**Total**: 16 distinct test scenarios, providing comprehensive performance profiles for both platforms.

### 3.6 Monitoring Tools and Metrics Collection

#### 3.6.1 Prometheus Metrics

**Scrape Configuration**:
- **Interval**: 5 seconds (high resolution for test observation)
- **Jobs**: Separate jobs for FastAPI and Golang services
- **Service discovery**: Uses Kubernetes service names (fastapi:80, golang:80)

**Key Metrics Collected**:

From **FastAPI** (via prometheus-fastapi-instrumentator):
- `http_requests_total{method, handler, status}` - Request count by endpoint
- `http_request_duration_seconds_bucket{le}` - Histogram of request durations
- `http_request_duration_seconds_sum` - Total time spent processing requests
- `http_request_duration_seconds_count` - Total request count
- `fastapi_inprogress{method, handler}` - Currently processing requests

From **Golang** (via custom middleware):
- `http_requests_total{method, endpoint, status}` - Request count by endpoint
- `http_request_duration_seconds_bucket{le}` - Histogram of request durations
- `http_request_duration_seconds_sum` - Total time spent
- `http_request_duration_seconds_count` - Total requests

**Histogram Buckets** (Prometheus default):
- 0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10 seconds
- Allows calculation of any percentile (p50, p90, p95, p99)

#### 3.6.2 Getting Valid Information from Prometheus

**During tests**:
1. **Request Rate**:
   ```promql
   rate(http_requests_total[1m])
   ```
   Shows requests per second over the last minute for each service and endpoint.

2. **Latency Percentiles**:
   ```promql
   histogram_quantile(0.90, rate(http_request_duration_seconds_bucket[1m]))
   histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[1m]))
   histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[1m]))
   ```
   Calculates p90, p95, p99 latencies from histogram data.

3. **Error Rate**:
   ```promql
   rate(http_requests_total{status=~"5.."}[1m]) / rate(http_requests_total[1m])
   ```
   Percentage of requests returning 5xx errors.

4. **Throughput Comparison**:
   ```promql
   sum(rate(http_requests_total{service="fastapi"}[1m])) vs
   sum(rate(http_requests_total{service="golang"}[1m]))
   ```
   Direct comparison of requests/second between services.

#### 3.6.3 k6 Output Metrics

k6 provides detailed test results including:

**From test result files** (example: `test-load-results/k6-golang-health-results.txt`):

1. **Execution Summary**:
   - VU progression through stages
   - Completed iterations
   - Test duration and timing

2. **HTTP Metrics**:
   - `http_req_duration`: avg, min, med, max, p(90), p(95)
   - `http_req_failed`: Percentage and count of failed requests
   - `http_reqs`: Total requests and requests per second
   - `data_sent` / `data_received`: Network throughput

3. **Connection Metrics**:
   - `http_req_blocked`: Time waiting for connection
   - `http_req_connecting`: TCP connection time
   - `http_req_sending` / `http_req_receiving`: Network transfer times
   - `http_req_waiting`: Time to first byte (TTFB)

4. **Threshold Validation**:
   - ✓ or ✗ indicators showing whether thresholds passed
   - Example: `✓ http_req_duration: avg=46.03ms ... p(90)=103.62ms p(95)=200.12ms`
   - Example: `✓ http_req_failed: 0.00%`

**Interpreting k6 results**:
- **Passed thresholds (✓)**: Service met performance expectations
- **Failed thresholds (✗)**: Service exceeded acceptable latency or error rates
- **Latency distribution**: Difference between median and p95/p99 shows consistency
- **Request rate**: Actual throughput achieved (e.g., 201.56 req/s)
- **VU behavior**: Compare actual iterations vs expected to identify blocking

#### 3.6.4 Grafana Dashboards

**Dashboard Configuration** (grafana-dashboard.json):
- Pre-configured panels for both services
- Time-series graphs showing:
  - Request rate over time
  - Latency percentiles (p50, p90, p95, p99) as separate lines
  - Error rate percentage
  - Request duration heatmaps
- Comparison views with FastAPI and Golang side-by-side
- Annotations for test start/end times

**During test observation**:
1. **Watch the request rate ramp**: Should follow the k6 stage pattern
2. **Monitor latency trends**: Look for degradation as VUs increase
3. **Check error rates**: Should remain at 0% for load tests
4. **Observe resource saturation**: CPU usage approaching 0.5 cores (100% of limit)

**Post-test analysis**:
1. **Compare median vs p95**: Large gaps indicate inconsistent performance
2. **Identify breaking points**: When did latency spike or errors appear?
3. **Resource correlation**: Did CPU hit 100% when latency degraded?
4. **Recovery behavior**: How quickly did latency return to normal during ramp-down?

### 3.7 Test Execution Workflow

1. **Setup**: Deploy services and monitoring stack to Kubernetes cluster
2. **Baseline**: Ensure services are healthy and Prometheus is scraping metrics
3. **Execute load tests**: Run k6 jobs for each service/endpoint combination
4. **Collect results**: k6 output saved to `test-load-results/` directory
5. **Execute stress tests**: Run k6 stress jobs
6. **Collect results**: k6 output saved to `test-stress-results/` directory
7. **Analysis**: Review k6 metrics, Prometheus data, and Grafana dashboards
8. **Comparison**: Generate comparative analysis between FastAPI and Golang across all scenarios

## 4. Expected Outcomes and Analysis Framework

This performance testing methodology provides empirical data to answer:

1. **Which platform handles baseline HTTP traffic more efficiently?** (/health endpoint)
2. **How does each platform perform under CPU-intensive workloads?** (/cpu endpoint)
3. **Which platform better handles concurrent I/O operations?** (/io endpoint)
4. **What are the JSON serialization performance differences?** (/json endpoint)
5. **Where are the breaking points under stress conditions?**
6. **How does performance scale with increasing concurrency?**
7. **What are the resource efficiency trade-offs?**

The comprehensive metrics collected enable data-driven decisions about backend platform selection based on:
- Application workload characteristics (CPU vs I/O bound)
- Performance requirements (latency SLOs, throughput needs)
- Resource constraints and cost considerations
- Operational complexity and deployment requirements

## 5. Conclusion

This research project establishes a rigorous, reproducible methodology for comparing FastAPI and Golang backend services under controlled conditions. By testing equivalent implementations across multiple workload types with identical resource constraints, the results provide actionable insights for technology selection in backend service development.

The combination of k6 load testing, Prometheus metrics collection, and Grafana visualization creates a comprehensive observability framework that captures both application-level performance and infrastructure-level resource utilization. This holistic approach ensures that performance decisions are based on empirical evidence rather than assumptions or anecdotal reports.

The findings from this research will inform architectural decisions for future backend service development, with clear performance profiles for each platform across different workload characteristics.

---

## References

- FastAPI Documentation: https://fastapi.tiangolo.com/
- Golang Documentation: https://go.dev/doc/
- k6 Documentation: https://k6.io/docs/
- Prometheus Documentation: https://prometheus.io/docs/
- Kubernetes Documentation: https://kubernetes.io/docs/
- Python asyncio: https://docs.python.org/3/library/asyncio.html
- Prometheus FastAPI Instrumentator: https://github.com/trallnag/prometheus-fastapi-instrumentator

## Appendices

### Appendix A: Test Result Locations
- Load test results: `/test-load-results/`
- Stress test results: `/test-stress-results/`
- Grafana dashboard: `/grafana-dashboard.json`

### Appendix B: Kubernetes Resources
- FastAPI deployment: `/k8s/fastapi-service/deployments.yaml`
- Golang deployment: `/k8s/go-service/deployments.yaml`
- k6 load jobs: `/k8s/k6-job.yaml`
- k6 stress jobs: `/k8s/k6-job-stress.yaml`
- Prometheus configuration: `/k8s/prometheus/prometheus-configmap.yaml`

### Appendix C: Application Code
- FastAPI implementation: `/apps/fastapi/app/main.py`
- Golang implementation: `/apps/golang/main.go`
