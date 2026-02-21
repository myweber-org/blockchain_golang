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
    tokenBytes := make([]byte, 16)
    if _, err := rand.Read(tokenBytes); err != nil {
        return "", err
    }
    token := hex.EncodeToString(tokenBytes)

    sessions[token] = Session{
        ID:        token,
        UserID:    userID,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }

    return token, nil
}

func ValidateSession(token string) (int, error) {
    session, exists := sessions[token]
    if !exists {
        return 0, errors.New("session not found")
    }

    if time.Now().After(session.ExpiresAt) {
        delete(sessions, token)
        return 0, errors.New("session expired")
    }

    return session.UserID, nil
}

func InvalidateSession(token string) {
    delete(sessions, token)
}