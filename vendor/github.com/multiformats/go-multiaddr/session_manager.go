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
	ErrInvalidToken    = errors.New("invalid token format")
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

func (m *Manager) generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (m *Manager) Create(userID string, data map[string]interface{}) (string, error) {
	token, err := m.generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		UserID:    userID,
		Data:      data,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.expiration),
	}

	ctx := context.Background()
	key := m.prefix + token

	err = m.client.Set(ctx, key, session, m.expiration).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *Manager) Get(token string) (*Session, error) {
	if len(token) < 10 {
		return nil, ErrInvalidToken
	}

	ctx := context.Background()
	key := m.prefix + token

	var session Session
	err := m.client.Get(ctx, key).Scan(&session)
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (m *Manager) Delete(token string) error {
	ctx := context.Background()
	key := m.prefix + token

	err := m.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Refresh(token string) error {
	session, err := m.Get(token)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(m.expiration)

	ctx := context.Background()
	key := m.prefix + token

	err = m.client.Set(ctx, key, session, m.expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Cleanup() error {
	ctx := context.Background()
	iter := m.client.Scan(ctx, 0, m.prefix+"*", 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		var session Session
		err := m.client.Get(ctx, key).Scan(&session)
		if err == nil && time.Now().After(session.ExpiresAt) {
			m.client.Del(ctx, key)
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}