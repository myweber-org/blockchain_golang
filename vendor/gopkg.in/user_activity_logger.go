package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	Method    string
	Path      string
	UserAgent string
	IPAddress string
	Duration  time.Duration
}

type ActivityLogger struct {
	activities chan ActivityLog
}

func NewActivityLogger(bufferSize int) *ActivityLogger {
	al := &ActivityLogger{
		activities: make(chan ActivityLog, bufferSize),
	}
	go al.processLogs()
	return al
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		next.ServeHTTP(w, r)
		
		activity := ActivityLog{
			Timestamp: time.Now(),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IPAddress: r.RemoteAddr,
			Duration:  time.Since(start),
		}
		
		select {
		case al.activities <- activity:
		default:
			log.Println("Activity log buffer full, dropping entry")
		}
	})
}

func (al *ActivityLogger) processLogs() {
	for activity := range al.activities {
		log.Printf("Activity: %s %s from %s (UA: %s) took %v",
			activity.Method,
			activity.Path,
			activity.IPAddress,
			activity.UserAgent,
			activity.Duration,
		)
	}
}

func (al *ActivityLogger) Close() {
	close(al.activities)
}