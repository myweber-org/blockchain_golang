package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type ActivityKey string

const UserActivityKey ActivityKey = "user_activity"

type UserActivity struct {
	UserID    string
	Action    string
	Timestamp time.Time
	IPAddress string
	UserAgent string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activity := UserActivity{
			UserID:    extractUserID(r),
			Action:    r.Method + " " + r.URL.Path,
			Timestamp: time.Now().UTC(),
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		ctx := context.WithValue(r.Context(), UserActivityKey, activity)
		next.ServeHTTP(w, r.WithContext(ctx))

		go logActivity(activity)
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return "authenticated_user"
	}
	return "anonymous"
}

func logActivity(activity UserActivity) {
	log.Printf("Activity: User=%s Action=%s IP=%s Time=%s",
		activity.UserID,
		activity.Action,
		activity.IPAddress,
		activity.Timestamp.Format(time.RFC3339))
}

func GetActivityFromContext(ctx context.Context) (UserActivity, bool) {
	activity, ok := ctx.Value(UserActivityKey).(UserActivity)
	return activity, ok
}package middleware

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
	StatusCode int
	Duration   time.Duration
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(lrw, r)
		
		duration := time.Since(start)
		activity := ActivityLog{
			Timestamp:  start,
			Method:     r.Method,
			Path:       r.URL.Path,
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			StatusCode: lrw.statusCode,
			Duration:   duration,
		}
		
		log.Printf("ACTIVITY: %s %s %d %s %s",
			activity.Method,
			activity.Path,
			activity.StatusCode,
			activity.Duration,
			activity.RemoteAddr,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}