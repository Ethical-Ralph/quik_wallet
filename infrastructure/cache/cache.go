package cache

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) string
}

func NewCache(cache Cache) Cache {
	return cache
}
