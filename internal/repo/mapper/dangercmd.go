package mapper

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"taskmanager/internal/models"
	"taskmanager/pkg/store"
)

var (
	defaultCommandMapper *DangerCmdMapper
	commandOnce          sync.Once
)

type DangerCmdMapper struct {
	BaseMapper
	mutex sync.Mutex
}

func (dm *DangerCmdMapper) Lock() {
	dm.mutex.Lock()
}

func (dm *DangerCmdMapper) Unlock() {
	dm.mutex.Unlock()
}

func NewCommandMapper() *DangerCmdMapper {
	dangerCmdMapper := &DangerCmdMapper{}
	dangerCmdMapper.BaseMapper.Lock = dangerCmdMapper.Lock
	dangerCmdMapper.BaseMapper.Unlock = dangerCmdMapper.Unlock
	return dangerCmdMapper
}

func GetDangerCmdMapper() *DangerCmdMapper {
	if defaultCommandMapper == nil {
		commandOnce.Do(func() {
			defaultCommandMapper = NewCommandMapper()
		})
	}
	return defaultCommandMapper
}

func (dm *DangerCmdMapper) Upsert(dangerCmd *models.DangerousCmd) error {
	if dangerCmd == nil {
		return errors.New("命令对象不能为空")
	}
	conflictKey := []clause.Column{
		{Name: "command"},
		{Name: "remarks"},
	}
	return dm.BaseMapper.Upsert(conflictKey, dangerCmd)
}

func (dm *DangerCmdMapper) Delete(filter *models.DangerousCmd) (deletedItems *[]models.DangerousCmd, err error) {
	if filter == nil {
		return deletedItems, errors.New("命令对象不能为空")
	}
	deletedItems = &[]models.DangerousCmd{}
	err = dm.BaseMapper.SoftDeleteByFilter(filter, deletedItems)
	return
}

func (dm *DangerCmdMapper) ListAllDangerousCommand(filter *models.DangerousCmd) ([]*models.DangerousCmd, error) {
	var result []*models.DangerousCmd
	_, err := dm.BaseMapper.FindAll(filter, &result)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}

func (dm *DangerCmdMapper) FindAllWithPager(filter *models.DangerousCmd, result *[]models.DangerousCmd, pageSize, pageNo int) error {
	_, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).
			Model(filter).
			Scopes(
				orderBy("createdAt"),
				paginate(pageSize, pageNo)).
			Find(result)
	})
	return err
}

func (dm *DangerCmdMapper) Count(filter *models.DangerousCmd, sortBy string) (int, error) {
	//count, err := dm.BaseMapper.Count(filter, sortBy)
	var count int64
	_, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Model(filter).
			Scopes(
				orderBy(sortBy),
			).Count(&count)
	})
	return int(count), err
}
