package main

import (
    "encoding/json"
    "sync"
    "time"
)

type UserPreferences struct {
    Theme     string `json:"theme"`
    Language  string `json:"language"`
    Timezone  string `json:"timezone"`
    NotificationsEnabled bool `json:"notifications_enabled"`
}

type PreferencesCache struct {
    mu        sync.RWMutex
    store     map[string]cachedPreferences
    ttl       time.Duration
}

type cachedPreferences struct {
    prefs      UserPreferences
    expiresAt  time.Time
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
    return &PreferencesCache{
        store: make(map[string]cachedPreferences),
        ttl:   ttl,
    }
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    cached, exists := c.store[userID]
    if !exists || time.Now().After(cached.expiresAt) {
        return UserPreferences{}, false
    }
    return cached.prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.store[userID] = cachedPreferences{
        prefs:     prefs,
        expiresAt: time.Now().Add(c.ttl),
    }
}

func (c *PreferencesCache) Invalidate(userID string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.store, userID)
}

func (c *PreferencesCache) Cleanup() {
    c.mu.Lock()
    defer c.mu.Unlock()

    now := time.Now()
    for userID, cached := range c.store {
        if now.After(cached.expiresAt) {
            delete(c.store, userID)
        }
    }
}

func main() {
    cache := NewPreferencesCache(5 * time.Minute)

    prefs := UserPreferences{
        Theme:     "dark",
        Language:  "en",
        Timezone:  "UTC",
        NotificationsEnabled: true,
    }

    cache.Set("user123", prefs)

    if cachedPrefs, found := cache.Get("user123"); found {
        data, _ := json.MarshalIndent(cachedPrefs, "", "  ")
        println("Cached preferences:", string(data))
    }

    cache.Invalidate("user123")
    cache.Cleanup()
}