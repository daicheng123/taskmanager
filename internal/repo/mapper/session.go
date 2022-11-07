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
	defaultSessionMapper *SessionMapper
	sessionOnce          sync.Once
)

type SessionMapper struct {
	BaseMapper
	mutex sync.Mutex
}

func (sm *SessionMapper) Lock() {
	sm.mutex.Lock()
}

func (sm *SessionMapper) Unlock() {
	sm.mutex.Unlock()
}

func NewSessionMapper() *SessionMapper {
	sessionMapper := &SessionMapper{}
	sessionMapper.BaseMapper.Lock = sessionMapper.Lock
	sessionMapper.BaseMapper.Unlock = sessionMapper.Unlock
	return sessionMapper
}

func GetSessionMapper() *SessionMapper {
	if defaultSessionMapper == nil {
		sessionOnce.Do(func() {
			defaultSessionMapper = NewSessionMapper()
		})
	}
	return defaultSessionMapper
}

func (sm *SessionMapper) FindOne(filter *models.SessionModel) (session *models.SessionModel, err error) {
	if filter == nil {
		filter = &models.SessionModel{}
	}
	session = &models.SessionModel{}
	_, err = sm.BaseMapper.FindOne(filter, session)
	return session, err
}

func (sm *SessionMapper) FindByToken(token string) (*models.SessionModel, error) {
	s := &models.SessionModel{
		SessionValue: token,
	}
	return sm.FindOne(s)
}

func (sm *SessionMapper) Save(session *models.SessionModel) error {
	if session == nil {
		return errors.New("session对象不能为空")
	}
	conflictKey := []clause.Column{
		//{Name: "user_id"},
		{Name: "session_value"},
		{Name: "expire_time"},
	}

	return sm.BaseMapper.Upsert(conflictKey, session)
}

// Delete 删除session
func (sm *SessionMapper) Delete(session *models.SessionModel) error {
	_, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Where(session).Delete(session)
	})
	return err
}
