package cache

import (
	"bytes"
	"encoding/gob"

	"go-web/internal/core/ports"

	"github.com/bradfitz/gomemcache/memcache"
)

type memCache struct {
	client *memcache.Client
}

func NewMemCache(addr string) ports.Cache {
	return &memCache{memcache.New(addr)}
}

func (c *memCache) Set(key string, value interface{}) error {
	if c == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return c.client.Set(&memcache.Item{Key: key, Value: buf.Bytes()})
}

func (c *memCache) Get(key string, value interface{}) error {
	if c == nil {
		return memcache.ErrCacheMiss
	}
	item, err := c.client.Get(key)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(item.Value)).Decode(value)
}

func (c *memCache) SetWithTTL(key string, value interface{}, ttl int) error {
	if c == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return c.client.Set(&memcache.Item{Key: key, Value: buf.Bytes(), Expiration: int32(ttl)})
}

func (c *memCache) Delete(key string) error {
	if c == nil {
		return nil
	}
	return c.client.Delete(key)
}
