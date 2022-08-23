package mapper

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"taskmanager/internal/models"
	"taskmanager/pkg/store"
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
	return user, err
}

//FindByEmail  通过email获取用户
func (um *UserMapper) FindByEmail(email string) (user *models.UserModel, err error) {
	filter := &models.UserModel{
		UserEmail: email,
	}
	return um.FindOne(filter)
}

//FindByUserId  通过用户id获取用户
func (um *UserMapper) FindByUserId(id uint) (user *models.UserModel, err error) {
	filter := &models.UserModel{
		BaseModel: models.BaseModel{
			ID: id,
		},
	}
	user = new(models.UserModel)
	_, err = um.BaseMapper.FindOne(filter, user)
	return
}

//FindByUserCode  通过user_code获取用户
func (um *UserMapper) FindByUserCode(code string) (user *models.UserModel, err error) {
	filter := &models.UserModel{
		UserCode: code,
	}
	return um.FindOne(filter)
}

//CreateUser  创建用户
func (um *UserMapper) CreateUser(user *models.UserModel) (err error) {
	if user == nil {
		return errors.New("用户对象不能为空")
	}
	_, err = store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Model(user).Create(user)
	})
	return err
}
