package worker

import (
	"context"
	"github.com/hibiken/asynq"
	"taskmanager/internal/conf"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/worker/tasks"
)

type Consumer struct {
	//queue queue.Queuer
	mux *asynq.ServeMux
	srv *asynq.Server
}

func NewConsumer() *Consumer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Network:  "tcp",
			Addr:     conf.GetRedisAddr(),
			Password: conf.GetRedisPasswd(),
			DB:       conf.GetTaskRedisDB(),
			PoolSize: 10,
		},
		asynq.Config{
			Concurrency: 10,
			IsFailure: func(err error) bool {
				if _, ok := err.(*tasks.RateLimitError); ok {
					return false
				}
				return true
			},
			Queues:         Queues,
			RetryDelayFunc: tasks.GetRetryDelay,
			Logger:         logger.GetAsynqLogger(),
		})

	mux := asynq.NewServeMux()
	c := &Consumer{
		mux: mux,
		srv: srv,
	}
	return c
}

func (c *Consumer) Start() {
	logger.Info("启动任务中心 worker...")
	if err := c.srv.Run(c.mux); err != nil {
		logger.Error("任务中心 worker 启动失败, err:[%s]", err.Error())
		return
	}
	return
}

func (c *Consumer) Stop() {
	c.srv.Stop()
	c.srv.Shutdown()
}

func (c *Consumer) RegisterHandlers(taskName TaskName, handler func(context.Context, *asynq.Task) error) {
	c.mux.HandleFunc(string(taskName), handler)
}
