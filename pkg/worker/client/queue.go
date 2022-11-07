package client

import (
	"encoding/json"
	"time"
)

type Queuer interface {
	Write(taskName string, queueName string, job *Job) error
	Options() QueueOptions
}

type Job struct {
	ID       string          `json:"id"`
	Payload  json.RawMessage `json:"payload"`
	Delay    time.Duration   `json:"delay"`
	MaxRetry int             `json:"max_retry"`
}

type QueueOptions struct {
	Names map[string]int
	Type  string
	//RedisClient       *rdb.Redis
	//RedisAddress      string
	//PrometheusAddress string
}
