package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserPreferences struct {
	UserID    string                 `json:"user_id"`
	Theme     string                 `json:"theme"`
	Language  string                 `json:"language"`
	Timezone  string                 `json:"timezone"`
	Settings  map[string]interface{} `json:"settings"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type PreferencesCache struct {
	redisClient *redis.Client
	db          *gorm.DB
	ttl         time.Duration
}

func NewPreferencesCache(redisClient *redis.Client, db *gorm.DB, ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		redisClient: redisClient,
		db:          db,
		ttl:         ttl,
	}
}

func (c *PreferencesCache) Get(ctx context.Context, userID string) (*UserPreferences, error) {
	cacheKey := "user_prefs:" + userID

	// Try cache first
	cachedData, err := c.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var prefs UserPreferences
		if err := json.Unmarshal([]byte(cachedData), &prefs); err == nil {
			return &prefs, nil
		}
	}

	// Cache miss, query database
	var prefs UserPreferences
	result := c.db.WithContext(ctx).Where("user_id = ?", userID).First(&prefs)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user preferences not found")
		}
		return nil, result.Error
	}

	// Update cache asynchronously
	go func() {
		ctx := context.Background()
		if data, err := json.Marshal(prefs); err == nil {
			c.redisClient.Set(ctx, cacheKey, data, c.ttl)
		}
	}()

	return &prefs, nil
}

func (c *PreferencesCache) Update(ctx context.Context, prefs *UserPreferences) error {
	// Update database
	prefs.UpdatedAt = time.Now()
	result := c.db.WithContext(ctx).Save(prefs)
	if result.Error != nil {
		return result.Error
	}

	// Update cache
	cacheKey := "user_prefs:" + prefs.UserID
	data, err := json.Marshal(prefs)
	if err != nil {
		return err
	}

	return c.redisClient.Set(ctx, cacheKey, data, c.ttl).Err()
}

func (c *PreferencesCache) Invalidate(ctx context.Context, userID string) error {
	cacheKey := "user_prefs:" + userID
	return c.redisClient.Del(ctx, cacheKey).Err()
}