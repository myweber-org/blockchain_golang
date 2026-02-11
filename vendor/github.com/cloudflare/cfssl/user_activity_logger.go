package middleware

import (
    "log"
    "net/http"
    "time"
)

type ActivityLogger struct {
    handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
    return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
    al.handler.ServeHTTP(rw, r)
    
    duration := time.Since(start)
    
    log.Printf("[%s] %s %s - %d - %s",
        time.Now().Format(time.RFC3339),
        r.Method,
        r.URL.Path,
        rw.statusCode,
        duration,
    )
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}