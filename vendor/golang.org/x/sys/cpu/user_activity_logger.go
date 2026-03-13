
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
	UserID    string
	IPAddress string
	UserAgent string
}

var activityChannel = make(chan ActivityLog, 100)

func init() {
	go processActivityLogs()
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserID(authHeader)
		}

		activity := ActivityLog{
			Timestamp: start,
			Method:    r.Method,
			Path:      r.URL.Path,
			UserID:    userID,
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		select {
		case activityChannel <- activity:
		default:
			log.Println("Activity log channel full, dropping log entry")
		}

		next.ServeHTTP(w, r)
	})
}

func extractUserID(token string) string {
	return "user_" + token[:8]
}

func processActivityLogs() {
	for activity := range activityChannel {
		log.Printf("ACTIVITY: %s %s %s %s %s\n",
			activity.Timestamp.Format(time.RFC3339),
			activity.Method,
			activity.Path,
			activity.UserID,
			activity.IPAddress,
		)
	}
}