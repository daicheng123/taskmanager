package mapper

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"taskmanager/internal/models"
	"taskmanager/pkg/store"
)

var (
	defaultTaskMapper *TaskMapper
	taskOnce          sync.Once
)

type TaskMapper struct {
	BaseMapper
	mutex sync.Mutex
}

func (tm *TaskMapper) Lock() {
	tm.mutex.Lock()
}

func (tm *TaskMapper) Unlock() {
	tm.mutex.Unlock()
}

func NewTaskMapper() *TaskMapper {
	taskMapper := &TaskMapper{}
	taskMapper.BaseMapper.Lock = taskMapper.Lock
	taskMapper.BaseMapper.Unlock = taskMapper.Unlock
	return taskMapper
}

func GetTaskMapper() *TaskMapper {
	if defaultTaskMapper == nil {
		taskOnce.Do(func() {
			defaultTaskMapper = NewTaskMapper()
		})
	}
	return defaultTaskMapper
}

func (tm *TaskMapper) CreateTask(value *models.Task) error {
	if value == nil {
		return errors.New("任务对象不能为空")
	}
	return tm.BaseMapper.Save(value)
}

func (tm *TaskMapper) CreateTaskStep(values *[]*models.TaskStep) error {
	if values == nil {
		return errors.New("任务步骤对象不能为空")
	}
	tm.Lock()
	defer tm.Unlock()
	_, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Model(&models.TaskStep{}).Create(values)
	})
	return err
}

func (tm *TaskMapper) PreLoadFindOne(ctx context.Context, filter *models.Task) (task *models.Task, err error) {
	if filter == nil {
		filter = &models.Task{}
	}
	task = &models.Task{}
	_, err = tm.BaseMapper.PreLoadFindOne(filter, task)
	if err != nil {
		return nil, err
	}
	return
}

func (tm *TaskMapper) UpdateTask(value *models.Task) error {
	if value == nil {
		return errors.New("任务对象不能为空")
	}
	_, err := tm.BaseMapper.Updates(value)
	return err
}

func (tm *TaskMapper) FindTaskStepOne(filter *models.TaskStep) (*models.TaskStep, error) {
	if filter == nil {
		filter = &models.TaskStep{}
	}
	result := &models.TaskStep{}
	_, err := tm.BaseMapper.FindOne(filter, result)
	return result, err
}

func (tm *TaskMapper) FindStepByExecutorAndTask(ipAddr string, taskName string) (*models.TaskStep, error) {
	ef := &models.Executor{
		IPAddr: ipAddr,
	}
	executor, err := GetExecutorMapper().FindOne(ef)
	if err != nil {
		fmt.Println("hahahahahahahahahahahahahahahahahahaha", err.Error())
		return nil, err
	}
	sf := &models.TaskStep{
		ExecutorID: executor.ID,
		TaskRefer:  taskName,
	}
	return tm.FindTaskStepOne(sf)
}

func (tm *TaskMapper) UpdateTaskStep(value models.UniqKeyGenerator) error {
	if value == nil {
		return errors.New("任务步骤对象不能为空")
	}
	_, err := tm.BaseMapper.Updates(value)
	return err
}
