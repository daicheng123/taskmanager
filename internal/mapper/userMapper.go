package mapper

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"taskmanager/internal/models"
)

var (
	defaultUserMapper *UserMapper
	userMapper        sync.Once
)

type UserMapper struct {
	BaseMapper
	mutex sync.Mutex
}

func (um *UserMapper) Lock() {
	um.mutex.Lock()
}

func (um *UserMapper) Unlock() {
	um.mutex.Unlock()
}

func NewUserMapper() *UserMapper {
	userMapper := &UserMapper{}
	userMapper.BaseMapper.Lock = userMapper.Lock
	userMapper.BaseMapper.Unlock = userMapper.Unlock
	return userMapper
}

func GetUserMapper() *UserMapper {
	if defaultUserMapper == nil {
		userMapper.Do(func() {
			defaultUserMapper = NewUserMapper()
		})
	}
	return defaultUserMapper
}

func (um *UserMapper) FindOne(filter *models.UserModel) (user *models.UserModel, err error) {
	if filter == nil {
		filter = &models.UserModel{}
	}
	user = &models.UserModel{}
	_, err = um.BaseMapper.FindOne(filter, user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	fmt.Println(err, "wowowo")
	return user, err
}

func (um *UserMapper) FindByEmail(email string) (user *models.UserModel, err error) {
	filter := &models.UserModel{
		UserEmail: email,
	}
	return um.FindOne(filter)
}
