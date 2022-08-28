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
	defaultTagsMapper *TagsMapper
	tagsOnce          sync.Once
)

type TagsMapper struct {
	BaseMapper
	mutex sync.Mutex
}

func (tm *TagsMapper) Lock() {
	tm.mutex.Lock()
}

func (tm *TagsMapper) Unlock() {
	tm.mutex.Unlock()
}

func NewTagsMapper() *TagsMapper {
	tagsMapper := &TagsMapper{}
	tagsMapper.BaseMapper.Lock = tagsMapper.Lock
	tagsMapper.BaseMapper.Unlock = tagsMapper.Unlock
	return tagsMapper
}

func GetTagsMapper() *TagsMapper {
	if defaultTagsMapper == nil {
		tagsOnce.Do(func() {
			defaultTagsMapper = NewTagsMapper()
		})
	}
	return defaultTagsMapper
}

func (tm *TagsMapper) Upsert(tag *models.TagsModel) error {
	if tag == nil {
		return errors.New("标签对象不能为空")
	}
	conflictKey := []clause.Column{
		{Name: "id"},
		{Name: "tag_name"},
	}
	return tm.BaseMapper.Upsert(conflictKey, tag)
}

func (tm *TagsMapper) FindAllWithPager(filter, result interface{}, pageSize, pageNo int,
	sortBy string, conditions, searches map[string]string) (*gorm.DB, error) {

	return store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).
			Model(filter).
			Scopes(
				conditionBy(conditions),
				searchBy(searches),
				orderBy(sortBy),
				paginate(pageSize, pageNo)).
			Find(result)
	})
}

func (tm *TagsMapper) Count(filter interface{}, sortBy string, conditions, searches map[string]string) (int, error) {
	var count int64
	_, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Model(filter).
			Scopes(
				orderBy(sortBy),
				conditionBy(conditions),
				searchBy(searches)).
			Count(&count)
	})
	return int(count), err
}

// Delete 软删除
func (tm *TagsMapper) Delete(filter *models.TagsModel) (*[]models.TagsModel, error) {
	deletedItems := &[]models.TagsModel{}
	if filter == nil {
		return deletedItems, nil
	}
	err := tm.BaseMapper.SoftDeleteByFilter(filter, deletedItems)
	return deletedItems, err
}
