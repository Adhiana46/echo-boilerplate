package cache

type Cache interface {
	Set(key string, value string, expSecond int32) error
	Get(key string) (string, error)
	Close() error
}
