package main

import (
	"log"
	"net/http"
	"time"
)

var (
	requestLatency = make(map[string]time.Duration)
	statusCodes    = make(map[int]int)
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		requestLatency[r.URL.Path] = duration
		statusCodes[recorder.statusCode]++

		log.Printf("Request to %s took %v, responded with %d", r.URL.Path, duration, recorder.statusCode)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, metrics!"))
	})
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("Slow endpoint"))
	})

	wrappedMux := metricsMiddleware(mux)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}