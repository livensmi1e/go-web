package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"go-web/internal/core/ports"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string, password string, db int) ports.Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &redisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (c *redisCache) Set(key string, value interface{}) error {
	if c == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, buf.Bytes(), 0).Err()
}

func (c *redisCache) Get(key string, value interface{}) error {
	if c == nil {
		return redis.Nil
	}
	data, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(value)
}

func (c *redisCache) SetWithTTL(key string, value interface{}, ttl int) error {
	if c == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, buf.Bytes(), time.Duration(ttl)*time.Second).Err()
}

func (c *redisCache) Delete(key string) error {
	if c == nil {
		return nil
	}
	return c.client.Del(c.ctx, key).Err()
}
