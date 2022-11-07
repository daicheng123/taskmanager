package client

import (
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"taskmanager/internal/conf"
	"taskmanager/pkg/logger"
)

var DefaultQueue Queuer

type RedisQueue struct {
	client    *asynq.Client
	inspector *asynq.Inspector
	opts      QueueOptions
}

func NewDefaultQueue(opts QueueOptions) {
	logger.Info("启动任务中心 queue")
	opt := &asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     conf.GetRedisAddr(),
		Password: conf.GetRedisPasswd(),
		DB:       conf.GetTaskRedisDB(),
		PoolSize: 10,
	}

	client := asynq.NewClient(opt)
	inspector := asynq.NewInspector(opt)

	DefaultQueue = &RedisQueue{
		client:    client,
		opts:      opts,
		inspector: inspector,
	}
}

func (q *RedisQueue) Write(taskName string, queueName string, job *Job) error {
	if job.ID == "" {
		job.ID = uuid.NewString()
	}

	t := asynq.NewTask(taskName, job.Payload, asynq.Queue(queueName), asynq.MaxRetry(job.MaxRetry), asynq.TaskID(job.ID), asynq.ProcessIn(job.Delay))
	_, err := q.client.Enqueue(t)
	return err
}

func (q *RedisQueue) StopTask(jobId string) error {
	return q.inspector.CancelProcessing(jobId)
}

func (q *RedisQueue) Options() QueueOptions {
	return q.opts
}
