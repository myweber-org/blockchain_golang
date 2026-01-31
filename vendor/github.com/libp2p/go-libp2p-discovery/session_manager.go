package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidToken    = errors.New("invalid session token")
)

type Session struct {
	UserID    string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Manager struct {
	client     *redis.Client
	expiration time.Duration
	prefix     string
}

func NewManager(client *redis.Client, expiration time.Duration) *Manager {
	return &Manager{
		client:     client,
		expiration: expiration,
		prefix:     "session:",
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (m *Manager) Create(userID, username string) (string, *Session, error) {
	token, err := generateToken()
	if err != nil {
		return "", nil, err
	}

	now := time.Now()
	session := &Session{
		UserID:    userID,
		Username:  username,
		CreatedAt: now,
		ExpiresAt: now.Add(m.expiration),
	}

	ctx := context.Background()
	key := m.prefix + token

	err = m.client.Set(ctx, key, userID, m.expiration).Err()
	if err != nil {
		return "", nil, fmt.Errorf("failed to store session: %w", err)
	}

	return token, session, nil
}

func (m *Manager) Validate(token string) (*Session, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	ctx := context.Background()
	key := m.prefix + token

	userID, err := m.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	ttl, err := m.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session TTL: %w", err)
	}

	now := time.Now()
	session := &Session{
		UserID:    userID,
		CreatedAt: now.Add(-m.expiration).Add(ttl),
		ExpiresAt: now.Add(ttl),
	}

	return session, nil
}

func (m *Manager) Refresh(token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	ctx := context.Background()
	key := m.prefix + token

	exists, err := m.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check session existence: %w", err)
	}

	if exists == 0 {
		return ErrSessionNotFound
	}

	err = m.client.Expire(ctx, key, m.expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	return nil
}

func (m *Manager) Delete(token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	ctx := context.Background()
	key := m.prefix + token

	_, err := m.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (m *Manager) Cleanup() error {
	ctx := context.Background()
	iter := m.client.Scan(ctx, 0, m.prefix+"*", 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		ttl, err := m.client.TTL(ctx, key).Result()
		if err != nil {
			continue
		}

		if ttl < 0 {
			m.client.Del(ctx, key)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed during cleanup scan: %w", err)
	}

	return nil
}