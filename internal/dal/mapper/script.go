package mapper

import (
	"errors"
	"sync"
	"taskmanager/internal/models"
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

func (sm *ScriptMapper) Save(value *models.Script) error {
	if value == nil {
		return errors.New("脚本对象不能为空")
	}
	return sm.BaseMapper.Save(value)
}
