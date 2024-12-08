package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetCache(rdClient *redis.Client, ctx context.Context, key string, ptr interface{}) error {
	cachedData, err := rdClient.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	return json.Unmarshal(cachedData, ptr)
}

func SetCache(rdClient *redis.Client, ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdClient.Set(ctx, key, jsonData, ttl).Err()
}

func ClearCacheByPattern(rdClient *redis.Client, ctx context.Context, pattern string) {
	keys, _ := rdClient.Keys(ctx, pattern).Result()
	for _, key := range keys {
		rdClient.Del(ctx, key)
	}
}
