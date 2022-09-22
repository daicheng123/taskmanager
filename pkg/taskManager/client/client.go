package main

import (
	"github.com/hibiken/asynq"
	"log"
	"taskmanager/pkg/taskManager/tasks"
	"time"
)

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "Dc!123",
		DB:       2,
		PoolSize: 10,
	})

	t1, err := tasks.NewWelcomeEmailTask(42)
	if err != nil {
		log.Fatal(err)
	}

	//t2, err := tasks.NewReminderEmailTask(42)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Process the task immediately.
	info, err := client.Enqueue(t1, asynq.MaxRetry(1), asynq.Timeout(time.Second*5))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("client执行欢迎邮件发送任务\n: %+v", info.Result)

	//info, err = client.Enqueue(t2, asynq.ProcessIn(2*time.Second))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("client执行重新发送邮件任务\n: %+v", info.Result)
}
