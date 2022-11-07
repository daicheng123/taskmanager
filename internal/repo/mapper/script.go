package mapper

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"sync"
	"taskmanager/internal/consts"
	"taskmanager/internal/models"
	"taskmanager/pkg/store"
)

var (
	defaultScriptMapper *ScriptMapper
	scriptOnce          sync.Once
)

type ScriptMapper struct {
	BaseMapper
	mutex sync.Mutex
}

func (sm *ScriptMapper) Lock() {
	sm.mutex.Lock()
}

func (sm *ScriptMapper) Unlock() {
	sm.mutex.Unlock()
}

func NewScriptMapper() *ScriptMapper {
	scriptMapper := &ScriptMapper{}
	scriptMapper.BaseMapper.Lock = scriptMapper.Lock
	scriptMapper.BaseMapper.Unlock = scriptMapper.Unlock
	return scriptMapper
}

func GetScriptMapper() *ScriptMapper {
	if defaultScriptMapper == nil {
		scriptOnce.Do(func() {
			defaultScriptMapper = NewScriptMapper()
		})
	}
	return defaultScriptMapper
}

func (sm *ScriptMapper) AddScript(value *models.Script) error {
	if value == nil {
		return errors.New("脚本对象不能为空")
	}

	sm.Lock()
	defer sm.Unlock()
	_, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Create(value)
	})
	return err
}

func (sm *ScriptMapper) DeleteAuditByFilter(audit *models.ScriptAudit) (err error) {
	if audit == nil {
		return errors.New("脚本审核对象不能为空")
	}
	audits := &[]models.ScriptAudit{}
	return sm.BaseMapper.SoftDeleteByFilter(audit, audits)
}

func (sm *ScriptMapper) UpdateScript(value *models.Script, auditor uint) (err error) {
	if value == nil {
		return errors.New("脚本审核对象不能为空")
	}

	filter := &models.ScriptAudit{ScriptRef: value.ID}
	audit, err := sm.FindAudit(filter)
	if err != nil {
		return err
	}
	switch value.Status {
	case consts.NoAudit:
		if audit != nil {
			audit.UserRef = auditor
			value.ScriptAudit = audit
		} else {
			value.ScriptAudit = &models.ScriptAudit{
				ScriptRef: value.ID,
				UserRef:   auditor,
				Reviewer:  value.LastOperator,
			}
		}
	case consts.PassAudit:
		if audit != nil {
			err = sm.DeleteAuditByFilter(audit)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("非法的审核状态")
	}
	_, err = sm.BaseMapper.Updates(value, "Tag")
	return err
}

func (sm *ScriptMapper) FindOne(filter *models.Script) (script *models.Script, err error) {
	if filter == nil {
		filter = &models.Script{}
	}
	script = &models.Script{}
	_, err = sm.BaseMapper.PreLoadFindOne(filter, script)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (sm *ScriptMapper) FindScriptOne(ctx context.Context, filter *models.Script) (script *models.Script, err error) {
	if filter == nil {
		filter = &models.Script{}
	}
	script = &models.Script{}
	_, err = sm.BaseMapper.PreLoadFindOne(filter, script)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (sm *ScriptMapper) FindAudit(filter *models.ScriptAudit) (audit *models.ScriptAudit, err error) {
	if filter == nil {
		filter = &models.ScriptAudit{}
	}
	audit = &models.ScriptAudit{}
	_, err = sm.BaseMapper.FindOne(filter, audit)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

// Delete 软删除
func (sm *ScriptMapper) Delete(filter *models.Script) (*[]models.Script, error) {
	deletedItems := &[]models.Script{}
	if filter == nil {
		return deletedItems, nil
	}
	err := sm.BaseMapper.SoftDeleteByFilter(filter, deletedItems)
	return deletedItems, err
}

func (sm *ScriptMapper) FindByName(name string) (script *models.Script, err error) {
	filter := &models.Script{
		ScriptName: name,
	}
	return sm.FindOne(filter)
}
