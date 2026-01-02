
package main

import (
    "encoding/json"
    "sync"
    "time"
)

type UserPreferences struct {
    UserID    string                 `json:"user_id"`
    Settings  map[string]interface{} `json:"settings"`
    UpdatedAt time.Time              `json:"updated_at"`
}

type PreferencesCache struct {
    mu      sync.RWMutex
    store   map[string]UserPreferences
    ttl     time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
    return &PreferencesCache{
        store: make(map[string]UserPreferences),
        ttl:   ttl,
    }
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    prefs, exists := c.store[userID]
    if !exists {
        return UserPreferences{}, false
    }

    if time.Since(prefs.UpdatedAt) > c.ttl {
        return UserPreferences{}, false
    }

    return prefs, true
}

func (c *PreferencesCache) Set(prefs UserPreferences) {
    c.mu.Lock()
    defer c.mu.Unlock()

    prefs.UpdatedAt = time.Now()
    c.store[prefs.UserID] = prefs
}

func (c *PreferencesCache) Invalidate(userID string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    delete(c.store, userID)
}

func (c *PreferencesCache) Size() int {
    c.mu.RLock()
    defer c.mu.RUnlock()

    return len(c.store)
}

func LoadPreferencesFromJSON(data []byte) (UserPreferences, error) {
    var prefs UserPreferences
    err := json.Unmarshal(data, &prefs)
    if err != nil {
        return UserPreferences{}, err
    }
    return prefs, nil
}