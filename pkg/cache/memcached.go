package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type mcCache struct {
	mc *memcache.Client
}

func NewMcCache(hosts ...string) (Cache, error) {
	mc := memcache.New(hosts...)

	return &mcCache{
		mc: mc,
	}, nil
}

func (c *mcCache) Set(key string, value string, expSecond int32) error {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Flags:      0,
		Expiration: expSecond,
	}

	return c.mc.Set(item)
}

func (c *mcCache) Get(key string) (string, error) {
	item, err := c.mc.Get(key)
	if err == nil {
		return string(item.Value), nil
	}

	if err == memcache.ErrCacheMiss {
		return "", ErrCacheNil
	}
	return "", err
}

func (c *mcCache) Delete(key string) error {
	err := c.mc.Delete(key)
	if err != nil {
		return err
	}

	return nil
}

func (c *mcCache) Close() error {
	return nil
}
