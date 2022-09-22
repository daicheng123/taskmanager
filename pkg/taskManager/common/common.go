package common

import (
	"github.com/hibiken/asynq"
	"taskmanager/internal/conf"
)

func NewClientOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     conf.GetRedisAddr(),
		Password: conf.GetRedisPasswd(),
		DB:       conf.GetTaskRedisDB(),
		//PoolSize: 10,
	}
}
