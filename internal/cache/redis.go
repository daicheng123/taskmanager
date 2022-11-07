package cache

import (
	"context"
	"sync"
	"taskmanager/internal/cache/utils"
	"taskmanager/pkg/store"
	"time"
)

type RedisCache struct {
	ctx  context.Context
	lock sync.RWMutex
}

func NewRedisCache(ctx context.Context) *RedisCache {
	return &RedisCache{
		ctx: ctx,
	}
}

func (r *RedisCache) Set(key string, value interface{}, attrs ...*utils.Attr) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	co, err := store.GetCacheOperator()
	if err != nil {
		return err
	}
	expire := utils.Attrs(attrs).Find(AttrExpire)
	if expire == nil {
		expire = 0
	}
	_, err = co.Set(r.ctx, key, value, expire.(time.Duration)).Result()
	return err
}

func (r *RedisCache) Get(key interface{}) (interface{}, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	co, err := store.GetCacheOperator()
	if err != nil {
		//if errors.Is(err, redis.Nil) {
		//
		//}
		return nil, err
	}
	value, err := co.Get(r.ctx, key.(string)).Result()
	if err != nil {
		return nil, err
	}
	return value, err
}

// Exists key是否存在
func (r *RedisCache) Exists(key string) (int64, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	co, err := store.GetCacheOperator()
	if err != nil {
		return 0, err
	}
	return co.Exists(r.ctx, key).Result()
}

func (r *RedisCache) HSet(key string, values ...interface{}) (int64, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	co, err := store.GetCacheOperator()
	if err != nil {
		return 0, err
	}
	return co.HSet(r.ctx, key, values).Result()
}

func (r *RedisCache) HGet(key string, field string) (string, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	co, err := store.GetCacheOperator()
	if err != nil {
		return "", err
	}
	return co.HGet(r.ctx, key, field).Result()
}
