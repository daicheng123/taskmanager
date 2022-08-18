package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"taskmanager/internal/conf"
	"time"
)

var cacheOperator *redis.Client
var once sync.Once

func NewCacheOperator() (err error) {
	once.Do(func() {
		options := &redis.Options{
			Network:            "tcp",
			Addr:               conf.GetRedisAddr(),
			DB:                 conf.GetRedisdb(),
			PoolSize:           15,
			MinIdleConns:       20,
			DialTimeout:        5 * time.Second,
			ReadTimeout:        3 * time.Second,
			WriteTimeout:       3 * time.Second,
			PoolTimeout:        4 * time.Second,
			IdleCheckFrequency: 60 * time.Second,
			IdleTimeout:        5 * time.Minute,
			MaxConnAge:         30 * time.Second,
			MaxRetries:         3,
			MinRetryBackoff:    8 * time.Millisecond,
			MaxRetryBackoff:    512 * time.Millisecond,
		}

		if conf.GetRedisUsePasswd() {
			options.Password = conf.GetRedisPasswd()
		}
		cacheOperator = redis.NewClient(options)

		_, err = cacheOperator.Ping(context.Background()).Result()
	})

	return err
}

func GetCacheOperator() (*redis.Client, error) {
	if cacheOperator == nil {
		if err := NewCacheOperator(); err != nil {
			return nil, fmt.Errorf("创建redis连接失败: %s", err.Error())
		}
	}
	return cacheOperator, nil
}
