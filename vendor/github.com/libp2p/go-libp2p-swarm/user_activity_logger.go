package middleware

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "github.com/go-redis/redis/v8"
    "golang.org/x/time/rate"
)

type ActivityLog struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
}

type ActivityLogger struct {
    redisClient *redis.Client
    limiter     *rate.Limiter
}

func NewActivityLogger(redisAddr string, rps int) *ActivityLogger {
    rdb := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: "",
        DB:       0,
    })

    return &ActivityLogger{
        redisClient: rdb,
        limiter:     rate.NewLimiter(rate.Limit(rps), rps),
    }
}

func (al *ActivityLogger) LogActivity(userID, action, ip, agent string) error {
    if !al.limiter.Allow() {
        return nil
    }

    logEntry := ActivityLog{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        IPAddress: ip,
        UserAgent: agent,
    }

    data, err := json.Marshal(logEntry)
    if err != nil {
        return err
    }

    ctx := context.Background()
    key := "activity:" + userID + ":" + time.Now().Format("20060102")
    return al.redisClient.RPush(ctx, key, data).Err()
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Header.Get("X-User-ID")
        if userID == "" {
            userID = "anonymous"
        }

        go func() {
            al.LogActivity(
                userID,
                r.Method+" "+r.URL.Path,
                r.RemoteAddr,
                r.UserAgent(),
            )
        }()

        next.ServeHTTP(w, r)
    })
}

func (al *ActivityLogger) GetRecentActivities(userID string, limit int64) ([]ActivityLog, error) {
    ctx := context.Background()
    key := "activity:" + userID + ":" + time.Now().Format("20060102")
    
    data, err := al.redisClient.LRange(ctx, key, -limit, -1).Result()
    if err != nil {
        return nil, err
    }

    activities := make([]ActivityLog, 0, len(data))
    for _, item := range data {
        var activity ActivityLog
        if err := json.Unmarshal([]byte(item), &activity); err == nil {
            activities = append(activities, activity)
        }
    }

    return activities, nil
}