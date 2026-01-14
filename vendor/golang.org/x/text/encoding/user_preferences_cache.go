package cache

import (
	"encoding/json"
	"time"
)

type UserPreferences struct {
	UserID      string
	Theme       string
	Language    string
	Timezone    string
	LastUpdated time.Time
}

type PreferencesCache struct {
	store map[string][]byte
	ttl   time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		store: make(map[string][]byte),
		ttl:   ttl,
	}
}

func (c *PreferencesCache) Get(userID string) (*UserPreferences, error) {
	data, exists := c.store[userID]
	if !exists {
		return nil, nil
	}

	var prefs UserPreferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		delete(c.store, userID)
		return nil, err
	}

	if time.Since(prefs.LastUpdated) > c.ttl {
		delete(c.store, userID)
		return nil, nil
	}

	return &prefs, nil
}

func (c *PreferencesCache) Set(prefs *UserPreferences) error {
	prefs.LastUpdated = time.Now()
	data, err := json.Marshal(prefs)
	if err != nil {
		return err
	}

	c.store[prefs.UserID] = data
	return nil
}

func (c *PreferencesCache) Invalidate(userID string) {
	delete(c.store, userID)
}

func (c *PreferencesCache) ClearExpired() {
	now := time.Now()
	for userID, data := range c.store {
		var prefs UserPreferences
		if err := json.Unmarshal(data, &prefs); err != nil {
			delete(c.store, userID)
			continue
		}
		if now.Sub(prefs.LastUpdated) > c.ttl {
			delete(c.store, userID)
		}
	}
}