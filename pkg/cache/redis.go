package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type redisCache struct {
	rdb *redis.Client
}

func NewRedisCache(host, port, password string, db int) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%v:%v", host, port),
		Password:    password,
		DB:          db,
		ReadTimeout: 0, // 3 second
	})

	return &redisCache{
		rdb: rdb,
	}, nil
}

func (c *redisCache) Set(key string, value string, expSecond int32) error {
	ctx := context.Background()

	return c.rdb.Set(ctx, key, value, time.Duration(expSecond)*time.Second).Err()
}

func (c *redisCache) Get(key string) (string, error) {
	ctx := context.Background()

	return c.rdb.Get(ctx, key).Result()
}

func (c *redisCache) Close() error {
	return c.rdb.Close()
}
