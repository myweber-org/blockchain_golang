package main

import (
    "encoding/json"
    "log"
    "os"
    "sync"
    "time"
)

type ActivityEvent struct {
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`
    Metadata  string    `json:"metadata,omitempty"`
}

type ActivityLogger struct {
    mu     sync.Mutex
    file   *os.File
    events []ActivityEvent
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &ActivityLogger{
        file:   file,
        events: make([]ActivityEvent, 0),
    }, nil
}

func (l *ActivityLogger) LogActivity(userID, eventType, metadata string) {
    l.mu.Lock()
    defer l.mu.Unlock()

    event := ActivityEvent{
        UserID:    userID,
        EventType: eventType,
        Timestamp: time.Now().UTC(),
        Metadata:  metadata,
    }

    l.events = append(l.events, event)

    data, err := json.Marshal(event)
    if err != nil {
        log.Printf("Failed to marshal event: %v", err)
        return
    }

    data = append(data, '\n')
    if _, err := l.file.Write(data); err != nil {
        log.Printf("Failed to write event: %v", err)
    }
}

func (l *ActivityLogger) GetRecentEvents(limit int) []ActivityEvent {
    l.mu.Lock()
    defer l.mu.Unlock()

    start := len(l.events) - limit
    if start < 0 {
        start = 0
    }
    return l.events[start:]
}

func (l *ActivityLogger) Close() error {
    return l.file.Close()
}

func main() {
    logger, err := NewActivityLogger("activity.log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    logger.LogActivity("user123", "login", "from web browser")
    logger.LogActivity("user123", "view_page", "/dashboard")
    logger.LogActivity("user456", "login", "from mobile app")

    recent := logger.GetRecentEvents(2)
    for _, event := range recent {
        log.Printf("Recent: %s - %s", event.UserID, event.EventType)
    }
}