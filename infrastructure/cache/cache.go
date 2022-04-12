package cache

type Cache interface {
	Set(key string, value interface{}) error
}

func NewCache(cache Cache) Cache {
	return cache
}
