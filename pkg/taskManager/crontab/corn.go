package main

import (
	"github.com/hibiken/asynq"

	"time"
)

// 从数据库加载定时任务
type configProvider struct {
}

func (p *configProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	var configs []*asynq.PeriodicTaskConfig
	return configs, nil
}

func InitPeriodicTask() error {
	//provider := &FileBasedConfigProvider{filename: "./periodic_task_config.yml"}
	provider := &configProvider{}
	mgr, err := asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			RedisConnOpt:               asynq.RedisClientOpt{Addr: "127.0.0.1:6379", Password: "Dc!123"},
			PeriodicTaskConfigProvider: provider,         // this provider object is the interface to your config source
			SyncInterval:               10 * time.Second, // this field specifies how often sync should happen
		})
	if err != nil {
		return err
	}

	if err := mgr.Run(); err != nil {
		return err
	}
	return nil
}

func ReloadPeriodicTask() error {
	return InitPeriodicTask()
}
