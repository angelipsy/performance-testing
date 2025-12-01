package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func prometheusMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(sw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(sw.statusCode)

		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	}
}

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func cpuHandler(w http.ResponseWriter, r *http.Request) {
	iterationsStr := r.URL.Query().Get("iterations")
	iterationsStr = map[bool]string{true: "100", false: iterationsStr}[iterationsStr == ""]
	iterations, _ := strconv.Atoi(iterationsStr)

	data := []byte("benchmark test data")
	for i := 0; i < iterations; i++ {
		h := sha256.Sum256(data)
		_ = hex.EncodeToString(h[:])
	}
	fmt.Fprint(w, "Completed "+strconv.Itoa(iterations)+" SHA256 hashes")
}

func ioHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(50 * time.Millisecond)
    fmt.Fprint(w, "ok")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	obj := map[string]interface{}{
		"name":  "test",
		"value": 123,
		"items": make([]int, 5000),
	}
	json.NewEncoder(w).Encode(obj)
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	const chunks = 20
	for i := 0; i < chunks; i++ {
		fmt.Fprintln(w, "chunk", i+1)
		flusher.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}



func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", prometheusMiddleware(healthHandler))
	mux.HandleFunc("/cpu", prometheusMiddleware(cpuHandler))
	mux.HandleFunc("/io", prometheusMiddleware(ioHandler))
	mux.HandleFunc("/json", prometheusMiddleware(jsonHandler))
	mux.HandleFunc("/stream", prometheusMiddleware(streamHandler))

	mux.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
