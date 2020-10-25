package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aidenwallis/customapi2/pkg/cache"
	"github.com/aidenwallis/customapi2/pkg/config"
	"github.com/go-redis/redis/v8"
)

type redisClient struct {
	client    *redis.Client
	keyPrefix string
}

func New(cfg *config.RedisConfig) cache.Cache {
	return &redisClient{
		client: redis.NewClient(&redis.Options{
			Addr:    cfg.Addr,
			Network: cfg.Network,
		}),
		keyPrefix: cfg.KeyPrefix,
	}
}

func (r *redisClient) Get(ctx context.Context, key string, out interface{}) error {
	res, err := r.client.Get(ctx, r.makeKey(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return cache.ErrNil
		}
		return err
	}
	if res == "null" {
		return cache.ErrNilResult
	}
	return json.Unmarshal([]byte(res), out)
}

func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.makeKey(key), string(bs), expiration).Err()
}

func (r *redisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.SetNX(ctx, r.makeKey(key), string(bs), expiration).Err()
}

func (r *redisClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.makeKey(key)).Err()
}

func (r *redisClient) makeKey(key string) string {
	return r.keyPrefix + key
}

func (r *redisClient) Close(ctx context.Context) error {
	return r.client.WithContext(ctx).Close()
}
