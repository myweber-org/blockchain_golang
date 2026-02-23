package session

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "time"

    "github.com/go-redis/redis/v8"
    "golang.org/x/net/context"
)

var (
    ErrInvalidToken = errors.New("invalid session token")
    ErrSessionExpired = errors.New("session expired")
)

type Session struct {
    UserID    string
    Username  string
    CreatedAt time.Time
    ExpiresAt time.Time
}

type Manager struct {
    client     *redis.Client
    prefix     string
    expiration time.Duration
}

func NewManager(client *redis.Client, prefix string, expiration time.Duration) *Manager {
    return &Manager{
        client:     client,
        prefix:     prefix,
        expiration: expiration,
    }
}

func generateToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func (m *Manager) Create(userID, username string) (string, error) {
    token, err := generateToken()
    if err != nil {
        return "", err
    }

    session := Session{
        UserID:    userID,
        Username:  username,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(m.expiration),
    }

    key := m.prefix + token
    ctx := context.Background()
    
    err = m.client.Set(ctx, key, session, m.expiration).Err()
    if err != nil {
        return "", err
    }

    return token, nil
}

func (m *Manager) Get(token string) (*Session, error) {
    key := m.prefix + token
    ctx := context.Background()

    var session Session
    err := m.client.Get(ctx, key).Scan(&session)
    if err != nil {
        if err == redis.Nil {
            return nil, ErrInvalidToken
        }
        return nil, err
    }

    if time.Now().After(session.ExpiresAt) {
        m.Delete(token)
        return nil, ErrSessionExpired
    }

    return &session, nil
}

func (m *Manager) Delete(token string) error {
    key := m.prefix + token
    ctx := context.Background()
    return m.client.Del(ctx, key).Err()
}

func (m *Manager) Refresh(token string) error {
    session, err := m.Get(token)
    if err != nil {
        return err
    }

    session.ExpiresAt = time.Now().Add(m.expiration)
    key := m.prefix + token
    ctx := context.Background()

    return m.client.Set(ctx, key, session, m.expiration).Err()
}