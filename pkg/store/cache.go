package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
	"sync"
	"taskmanager/internal/conf"
	"time"
)

var cacheOperator redis.UniversalClient
var once sync.Once

func NewCacheOperator() (err error) {
	once.Do(func() {
		if conf.GetRedisIsCluster() {
			options := &redis.ClusterOptions{
				//Network:            "tcp",
				Addrs: strings.Split(conf.GetRedisAddr(), ","),
				//:                 conf.GetRedisdb(),
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

			cacheOperator = redis.NewClusterClient(options)
		} else {
			options := &redis.Options{
				Network:            "tcp",
				Addr:               strings.Split(conf.GetRedisAddr(), ",")[0],
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
			cacheOperator = redis.NewClient(options)
		}

		_, err = cacheOperator.Ping(context.Background()).Result()
	})

	return err
}

func GetCacheOperator() (redis.UniversalClient, error) {
	if cacheOperator == nil {
		if err := NewCacheOperator(); err != nil {
			return nil, fmt.Errorf("创建redis连接失败: %s", err.Error())
		}
	}
	return cacheOperator, nil
}
