package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp  time.Time
	Method     string
	Path       string
	RemoteAddr string
	UserAgent  string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activity := ActivityLog{
			Timestamp:  time.Now().UTC(),
			Method:     r.Method,
			Path:       r.URL.Path,
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
		}

		log.Printf("Activity: %s %s from %s (UA: %s) at %s",
			activity.Method,
			activity.Path,
			activity.RemoteAddr,
			activity.UserAgent,
			activity.Timestamp.Format(time.RFC3339),
		)

		next.ServeHTTP(w, r)
	})
}