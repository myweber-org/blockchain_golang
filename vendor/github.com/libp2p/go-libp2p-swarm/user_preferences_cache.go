package main

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type UserPreferences struct {
	UserID    string
	Theme     string
	Language  string
	Timezone  string
	UpdatedAt time.Time
}

type PreferencesCache struct {
	mu      sync.RWMutex
	store   map[string]UserPreferences
	ttl     time.Duration
	cleanup *time.Ticker
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	cache := &PreferencesCache{
		store: make(map[string]UserPreferences),
		ttl:   ttl,
	}
	if ttl > 0 {
		cache.cleanup = time.NewTicker(ttl * 2)
		go cache.startCleanup()
	}
	return cache
}

func (c *PreferencesCache) Set(key string, prefs UserPreferences) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if key == "" {
		return errors.New("empty key not allowed")
	}

	prefs.UpdatedAt = time.Now()
	c.store[key] = prefs
	return nil
}

func (c *PreferencesCache) Get(key string) (UserPreferences, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	prefs, exists := c.store[key]
	if !exists {
		return UserPreferences{}, false
	}

	if c.ttl > 0 && time.Since(prefs.UpdatedAt) > c.ttl {
		return UserPreferences{}, false
	}

	return prefs, true
}

func (c *PreferencesCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func (c *PreferencesCache) startCleanup() {
	for range c.cleanup.C {
		c.mu.Lock()
		now := time.Now()
		for key, prefs := range c.store {
			if now.Sub(prefs.UpdatedAt) > c.ttl {
				delete(c.store, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *PreferencesCache) MarshalJSON() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return json.Marshal(c.store)
}

func (c *PreferencesCache) Close() {
	if c.cleanup != nil {
		c.cleanup.Stop()
	}
}