package mapper

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"taskmanager/internal/models"
	"taskmanager/pkg/store"
)

var (
	defaultExecutorMapper *ExecutorMapper
	execOnce              sync.Once
)

type ExecutorMapper struct {
	BaseMapper
	mutex sync.Mutex
}

func (em *ExecutorMapper) Lock() {
	em.mutex.Lock()
}

func (em *ExecutorMapper) Unlock() {
	em.mutex.Unlock()
}

func NewExecutorMapper() *ExecutorMapper {
	executorMapper := &ExecutorMapper{}
	executorMapper.BaseMapper.Lock = executorMapper.Lock
	executorMapper.BaseMapper.Unlock = executorMapper.Unlock
	return executorMapper
}

func GetExecutorMapper() *ExecutorMapper {
	if defaultExecutorMapper == nil {
		execOnce.Do(func() {
			defaultExecutorMapper = NewExecutorMapper()
		})
	}
	return defaultExecutorMapper
}

func (em *ExecutorMapper) Upsert(executor *models.Executor) error {
	if executor == nil {
		return errors.New("执行器对象不能为空")
	}
	return em.BaseMapper.Save(executor)
}

func (em *ExecutorMapper) FindOne(filter *models.Executor) (executor *models.Executor, err error) {
	if filter == nil {
		filter = &models.Executor{}
	}
	executor = &models.Executor{}
	_, err = em.BaseMapper.FindOne(filter, executor)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return executor, err
}

func (em *ExecutorMapper) PreLoadFindOne(filter *models.Executor) (executor *models.Executor, err error) {
	if filter == nil {
		filter = &models.Executor{}
	}
	executor = &models.Executor{}
	_, err = em.BaseMapper.PreLoadFindOne(filter, executor)
	//if errors.Is(err, gorm.ErrRecordNotFound) {
	//	return nil, nil
	//}
	return
}

func (em *ExecutorMapper) Updates(value *models.Executor) error {
	_, err := em.BaseMapper.Updates(value)
	return err
}

func (em *ExecutorMapper) FindAllWithPager(filter, result interface{}, pageSize, pageNo int,
	sortBy string, conditions, searches map[string]interface{}) error {
	_, err := em.BaseMapper.FindAllWithPager(filter, result, pageSize, pageNo,
		sortBy, conditions, searches)
	return err
}

func (em *ExecutorMapper) Count(filter interface{}, sortBy string, conditions, searches map[string]interface{}) (int, error) {
	count, err := em.BaseMapper.Count(filter, sortBy, conditions, searches)
	return count, err
}

// Delete 软删除
func (em *ExecutorMapper) Delete(filter *models.Executor) (*[]models.Executor, error) {
	deletedItems := &[]models.Executor{}
	if filter == nil {
		return deletedItems, nil
	}
	err := em.BaseMapper.SoftDeleteByFilter(filter, deletedItems)
	return deletedItems, err
}

// BatchDeleteById 批量软删除
func (em *ExecutorMapper) BatchDeleteById(ids ...uint) error {
	err := store.Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			filter := &models.Executor{
				BaseModel: &models.BaseModel{ID: id}}
			if _, err := em.Delete(filter); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
