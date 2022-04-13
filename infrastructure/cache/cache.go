package cache

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) string
	Del(key string)
}

func NewCache(cache Cache) Cache {
	return cache
}
