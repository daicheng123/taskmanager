package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	mapper2 "taskmanager/internal/dal/mapper"
	"taskmanager/internal/models"
	"taskmanager/internal/service"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
	"taskmanager/utils"
	"time"
)

const (
	UserTokenStr = "USER_TOKEN"
)

type UserLoginService struct {
	Password string `json:"userPassword"  binding:"required"`
	Email    string `json:"email"         binding:"required"`
}

func NewUserLoginService() *UserLoginService {
	return &UserLoginService{}
}

func (ul *UserLoginService) Login() *serializer.Response {
	user, err := mapper2.GetUserMapper().FindByEmail(ul.Email)
	if err != nil {
		return serializer.DBErr("查询用户失败", err)
	}

	if user == nil {
		return serializer.Err(serializer.CodeParamErr, "用户不存在", nil)
	}
	originPass, err := utils.Decrypt(user.UserPassword)
	if err != nil {
		logger.Error("解析密码错误, err:[%s]", err.Error())
		return serializer.Err(serializer.CodeServerInternalError, "解析密码错误", err)
	}
	if originPass != ul.Password {
		return serializer.Err(serializer.CodeParamErr, "错误的用户密码", nil)
	}

	// 登录
	sv := utils.GenMd5(user.UserCode)
	session := &models.SessionModel{
		UserID:       user.ID,
		SessionValue: sv,
		ExpireTime:   time.Now().Add(service.SessionCookieAge * time.Second),
	}
	err = mapper2.GetSessionMapper().Save(session)
	if err != nil {
		logger.Error("创建会话失败,err:[%s]", err.Error())
		return serializer.DBErr("创建会话失败", err)
	}
	return &serializer.Response{Data: sv}
}

func (ul *UserLoginService) UserInfo(ctx *gin.Context) *serializer.Response {
	token, exist := ctx.Get(UserTokenStr)
	fmt.Println(token, exist, "hahahahahaha")
	if !exist || !service.SessionJudge(token.(string)) {
		return serializer.Err(serializer.CodeCheckLogin, "用户未登录", nil)
	}
	session, err := mapper2.NewSessionMapper().FindByToken(token.(string))
	if err != nil {
		logger.Error("用户session查询失败, err:[%s]", err.Error())
		return serializer.DBErr("用户信息查询失败", err)
	}
	user, err := mapper2.NewUserMapper().FindByUserId(session.UserID)
	if err != nil {
		logger.Error("用户信息查询失败, err:[%s]", err.Error())
		return serializer.DBErr("用户信息查询失败", err)
	}
	return &serializer.Response{Data: user}
}

func (ul *UserLoginService) UserLogOut(token string) *serializer.Response {
	session := &models.SessionModel{
		SessionValue: token,
	}
	err := mapper2.GetSessionMapper().Delete(session)
	if err != nil {
		logger.Error("用户session删除失败, err:[%s]", err.Error())
		return serializer.DBErr("用户session清理失败", err)
	}
	return &serializer.Response{}
}
