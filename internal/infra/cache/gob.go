package cache

import (
	"bytes"
	"encoding/gob"

	"go-web/internal/core/ports"

	"github.com/bradfitz/gomemcache/memcache"
)

type gobCache struct {
	client *memcache.Client
}

func NewGobCache(addr string) ports.Cache {
	return &gobCache{memcache.New(addr)}
}

func (c *gobCache) Set(key string, value interface{}) error {
	if c == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return c.client.Set(&memcache.Item{Key: key, Value: buf.Bytes()})
}

func (c *gobCache) Get(key string, value interface{}) error {
	if c == nil {
		return memcache.ErrCacheMiss
	}
	item, err := c.client.Get(key)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(item.Value)).Decode(value)
}
