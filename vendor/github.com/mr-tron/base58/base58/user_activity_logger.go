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
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activity := UserActivity{
			UserID:    extractUserID(r),
			Action:    r.Method + " " + r.URL.Path,
			Timestamp: time.Now().UTC(),
			IPAddress: r.RemoteAddr,
		}

		ctx := context.WithValue(r.Context(), UserActivityKey, activity)
		next.ServeHTTP(w, r.WithContext(ctx))

		logActivity(activity)
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
		activity.Timestamp.Format(time.RFC3339),
	)
}

func GetActivityFromContext(ctx context.Context) (UserActivity, bool) {
	activity, ok := ctx.Value(UserActivityKey).(UserActivity)
	return activity, ok
}