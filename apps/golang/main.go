package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
	"path/filepath"
	"io/ioutil"
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

type Preferences struct {
	Theme         string `json:"theme"`
	Notifications bool   `json:"notifications"`
}

type Metadata struct {
	Role        string      `json:"role"`
	CreatedAt   string      `json:"created_at"`
	Preferences Preferences `json:"preferences"`
}

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Active   bool     `json:"active"`
	Metadata Metadata `json:"metadata"`
}

type Pagination struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type JSONResponse struct {
	Status     string     `json:"status"`
	Timestamp  string     `json:"timestamp"`
	Users      []User     `json:"users"`
	Pagination Pagination `json:"pagination"`
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
	// Write data to a temporary file
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, "io_test.txt")

	// Write 1000 lines to the file (line-by-line)
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(file, "Line %d: This is test data for I/O operations\n", i)
	}
	file.Close()

	// Read the file back to simulate I/O
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Clean up
	os.Remove(filePath)

	// Count lines
	lines := 0
	for _, b := range data {
		if b == '\n' {
			lines++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":        "ok",
		"lines_written": lines,
	}
	json.NewEncoder(w).Encode(response)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	// Create a complex nested data structure
	users := make([]User, 1000)
	for i := 0; i < 1000; i++ {
		role := "user"
		if i%10 == 0 {
			role = "admin"
		}

		users[i] = User{
			ID:     i,
			Name:   fmt.Sprintf("User %d", i),
			Email:  fmt.Sprintf("user%d@example.com", i),
			Active: i%2 == 0,
			Metadata: Metadata{
				Role:      role,
				CreatedAt: "2024-01-01",
				Preferences: Preferences{
					Theme:         "dark",
					Notifications: true,
				},
			},
		}
	}

	data := JSONResponse{
		Status:    "success",
		Timestamp: "2024-01-01T00:00:00Z",
		Users:     users,
		Pagination: Pagination{
			Total:   1000,
			Page:    1,
			PerPage: 1000,
		},
	}

	// Explicitly serialize to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to serialize JSON", http.StatusInternalServerError)
		return
	}

	// Deserialize to demonstrate serialization (similar to Python implementation)
	var result JSONResponse
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		http.Error(w, "Failed to deserialize JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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
