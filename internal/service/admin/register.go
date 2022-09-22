package admin

import (
	"gopkg.in/errgo.v2/errors"
	"taskmanager/internal/dal/mapper"
	"taskmanager/internal/models"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
	"taskmanager/utils"
)

type UserRegisterService struct {
	*models.UserModel `json:",inline"  binding:"required"`
	Password          string `json:"userPassword"  binding:"required"`
	EmailCode         string `json:"emailCode"     binding:"required"`
}

func NewUserRegisterService() *UserRegisterService {
	return &UserRegisterService{}
}

func (us *UserRegisterService) GetUserByUserCode(code string) *serializer.Response {
	user, err := mapper.GetUserMapper().FindByEmail(code)
	if err != nil {
		return serializer.DBErr("查询用户失败", err)
	}
	return &serializer.Response{Data: user}
}

func (us *UserRegisterService) AddUser() *serializer.Response {
	var (
		err = errors.New("")
	)

	if us.Password != "" {
		us.UserModel.UserPassword, err = utils.Encrypt(us.Password)
		if err != nil {
			return serializer.Err(serializer.CodeServerInternalError, "用户密码设置失败", err)
		}
	}
	if err = mapper.GetUserMapper().CreateUser(us.UserModel); err != nil {
		logger.Error("创建用户失败,error:[%s]", err.Error())
		return serializer.DBErr("创建用户失败", err)
	}
	return &serializer.Response{}
}
