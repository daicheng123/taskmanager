package worker

import (
	"taskmanager/pkg/worker/client"
	"taskmanager/pkg/worker/tasks"
	"taskmanager/utils"
)

func InitWorker() {
	queueOpts := client.QueueOptions{
		Names: Queues,
		Type:  RedisQueueProvider,
	}

	go utils.RunSafeWithMsg(func() {
		client.NewDefaultQueue(queueOpts)
	}, "任务中心 queue启动失败")

	go utils.RunSafeWithMsg(func() {
		consumer := NewConsumer()
		consumer.RegisterHandlers(RunAnsibleAdhoc, tasks.RunAnsibleAdhoc())
		consumer.Start()
	}, "任务中心 worker 启动失败")

}
