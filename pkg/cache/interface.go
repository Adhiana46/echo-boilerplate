package cache

type Cache interface {
	Set(key string, value any) error
	Get(key string) (any, error)
}
