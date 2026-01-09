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
    ErrSessionNotFound = errors.New("session not found")
    ErrInvalidToken    = errors.New("invalid session token")
)

type Session struct {
    UserID    string                 `json:"user_id"`
    Data      map[string]interface{} `json:"data"`
    CreatedAt time.Time              `json:"created_at"`
    ExpiresAt time.Time              `json:"expires_at"`
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

func (m *Manager) Create(userID string, data map[string]interface{}) (string, error) {
    token, err := generateToken()
    if err != nil {
        return "", err
    }

    session := Session{
        UserID:    userID,
        Data:      data,
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
    if token == "" {
        return nil, ErrInvalidToken
    }

    key := m.prefix + token
    ctx := context.Background()

    var session Session
    err := m.client.Get(ctx, key).Scan(&session)
    if err != nil {
        if err == redis.Nil {
            return nil, ErrSessionNotFound
        }
        return nil, err
    }

    if time.Now().After(session.ExpiresAt) {
        m.Delete(token)
        return nil, ErrSessionNotFound
    }

    return &session, nil
}

func (m *Manager) Delete(token string) error {
    if token == "" {
        return ErrInvalidToken
    }

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

func (m *Manager) Cleanup() error {
    ctx := context.Background()
    return m.client.FlushDB(ctx).Err()
}