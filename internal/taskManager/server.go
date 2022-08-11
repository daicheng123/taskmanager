package taskManager

import (
	"github.com/RichardKnop/machinery/v2"
	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	"github.com/RichardKnop/machinery/v2/config"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/RichardKnop/machinery/v2/tasks"
	"sync"
	"taskmanager/internal/conf"
	"taskmanager/pkg/logger"
	"taskmanager/utils"
)

var (
	defaultTaskCenter *taskCenter
	taskLock          sync.Mutex
	//asyncTaskMap      map[string]interface{}
)

type taskCenter struct {
	server *machinery.Server
	worker *machinery.Worker
}

func GetDefaultTaskCenter() *taskCenter {
	if defaultTaskCenter == nil {
		taskLock.Lock()
		defer taskLock.Unlock()
		if defaultTaskCenter == nil {
			defaultTaskCenter = newTaskCenter()
		}
	}
	return defaultTaskCenter
}

func newTaskCenter() *taskCenter {
	var (
		redisAddr   = conf.GetRedisAddr()
		redisPasswd = conf.GetRedisPasswd()
		redisDB     = conf.GetRedisdb()
	)
	addr := redisPasswd + "@" + redisAddr

	tc := new(taskCenter)
	cf := tc.config()
	if cf == nil {
		return nil
	}

	broker := redisbroker.NewGR(cf, []string{addr}, redisDB)
	backend := redisbackend.NewGR(cf, []string{addr}, redisDB)
	lock := eagerlock.New()

	tc.server = machinery.NewServer(cf, broker, backend, lock)
	// 初始化任务
	utils.ShouldShutDown(tc.registerTasks())
	// 初始化周期任务
	utils.ShouldShutDown(tc.applyPeriodicTask())
	return tc
}

func (tc *taskCenter) config() *config.Config {
	return &config.Config{
		DefaultQueue:    "task_manager",
		ResultsExpireIn: 3600,
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}
}

func (tc *taskCenter) registerTasks() error {
	//logger.TaskInfo("register tasks")
	//return tc.server.RegisterTask("add", Add)
	return nil
}

func (tc *taskCenter) applyPeriodicTask() error {
	return nil
}

func (tc *taskCenter) StartWorker(concurrency int) {
	consumerTag := "taskWorker"
	tc.worker = tc.server.NewWorker(consumerTag, concurrency)

	tc.worker.SetErrorHandler(func(err error) {
		logger.TaskError("执行失败: %s", err.Error())
	})

	tc.worker.SetPostTaskHandler(func(signature *tasks.Signature) {
		logger.TaskWarning("执行结束: %s", signature.Name)
	})

	tc.worker.SetPreTaskHandler(func(signature *tasks.Signature) {
		logger.TaskInfo("开始执行: %s", signature.Name)
	})

	utils.ShouldShutDown(tc.worker.Launch())
}
