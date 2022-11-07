package worker

type QueueName string

const (
	AsyncQueue         QueueName = "AsyncQueue"
	ScheduleQueue      QueueName = "ScheduleQueue"
	DefaultQueue       QueueName = "DefaultQueue"
	RedisQueueProvider string    = "redis"
)

var Queues = map[string]int{
	string(AsyncQueue):    5,
	string(ScheduleQueue): 3,
	string(DefaultQueue):  2,
}

type TaskName string

const (
	RunAnsibleAdhoc TaskName = "runAnsibleAdhocProcess"
)
