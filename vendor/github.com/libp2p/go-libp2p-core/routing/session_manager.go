package session

import (
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"
)

type Session struct {
    ID        string
    UserID    int
    ExpiresAt time.Time
}

var sessions = make(map[string]Session)

func GenerateSession(userID int) (string, error) {
    token := make([]byte, 16)
    _, err := rand.Read(token)
    if err != nil {
        return "", err
    }

    sessionID := hex.EncodeToString(token)
    expiresAt := time.Now().Add(24 * time.Hour)

    sessions[sessionID] = Session{
        ID:        sessionID,
        UserID:    userID,
        ExpiresAt: expiresAt,
    }

    return sessionID, nil
}

func ValidateSession(sessionID string) (Session, error) {
    session, exists := sessions[sessionID]
    if !exists {
        return Session{}, errors.New("session not found")
    }

    if time.Now().After(session.ExpiresAt) {
        delete(sessions, sessionID)
        return Session{}, errors.New("session expired")
    }

    return session, nil
}

func InvalidateSession(sessionID string) {
    delete(sessions, sessionID)
}

func CleanupExpiredSessions() {
    now := time.Now()
    for id, session := range sessions {
        if now.After(session.ExpiresAt) {
            delete(sessions, id)
        }
    }
}