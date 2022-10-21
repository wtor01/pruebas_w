package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"time"
)

type DataCacheRedis struct {
	r *redis.Client
}

func NewDataCacheRedis(r *redis.Client) *DataCacheRedis {
	if r == nil {
		return nil
	}
	return &DataCacheRedis{r: r}
}

func (rc *DataCacheRedis) Get(ctx context.Context, valueParam any, keySearch string) error {
	value, err := rc.r.Get(ctx, keySearch).Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(value, valueParam)
	return err
}

func (rc *DataCacheRedis) Clean(ctx context.Context, keySearch string) error {
	keys, _, err := rc.r.Scan(ctx, 0, keySearch, 0).Result()
	if err != nil {
		return err
	}
	rc.r.Del(ctx, keys...)
	return nil
}

func (rc *DataCacheRedis) Set(ctx context.Context, keySearch string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	rc.r.Set(ctx, keySearch, data, 0)
	return nil
}

func (rc *DataCacheRedis) SetWithExpiration(ctx context.Context, keySearch string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	rc.r.Set(ctx, keySearch, data, expiration)
	return nil
}
