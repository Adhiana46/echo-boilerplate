package cache

import "errors"

var (
	ErrCacheNil = errors.New("Nil")
)

type Cache interface {
	Set(key string, value string, expSecond int32) error
	Get(key string) (string, error)
	Delete(key string) error
	Close() error
}
