package ports

type Cache interface {
	Set(key string, value interface{}) error
	SetWithTTL(key string, value interface{}, ttl int) error
	Get(key string, value interface{}) error
}
