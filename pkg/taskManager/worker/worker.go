package main

import (
	"github.com/hibiken/asynq"
	"log"
	"taskmanager/pkg/taskManager/tasks"
	"time"
)

func main() {
	srv := asynq.NewServer(asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "Dc!123",
		DB:       2,
		PoolSize: 10,
	}, asynq.Config{
		Concurrency:              10,
		DelayedTaskCheckInterval: 1,
		RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
			return asynq.DefaultRetryDelayFunc(n, e, t)
		}})
	mux := asynq.NewServeMux()

	mux.Handle(tasks.TypeWelcomeEmail, tasks.NewEmailProcessor())

	//mux.HandleFunc(tasks.TypeReminderEmail, tasks.HandleReminderEmailTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
