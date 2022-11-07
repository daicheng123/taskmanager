package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/cache"
	"taskmanager/internal/service/admin"
	"taskmanager/internal/web"
	webutils "taskmanager/internal/web/utils"
	"taskmanager/pkg/serializer"
	"taskmanager/utils"
)

const (
	UserControllerGroup = "users"
)

type UserController struct {
	service *admin.UserRegisterService
}

func NewUserController() *UserController {
	return &UserController{
		service: new(admin.UserRegisterService),
	}
}

func (uc *UserController) checkUserExist(ctx *gin.Context) {
	service := admin.NewUserRegisterService()
	ctx.JSON(http.StatusOK, service.GetUserByUserCode(ctx.Param("userCode")))
}

// UserRegister 用户注册
func (uc *UserController) userRegister(ctx *gin.Context) {
	service := admin.NewUserRegisterService()
	err := ctx.ShouldBindJSON(service)
	if err != nil {
		ctx.JSON(http.StatusOK, webutils.ErrorResponse(err))
	}
	// 校验验证码是否过期
	//codeKey := service.UserEmail + ""
	//r := cache.NewStringOperation().Get(codeKey).UnwrapOrElse(func(err error) {
	//	logger.Error("获取邮箱验证码失败: err:[%s]", err.Error())
	//})

	// 校验验证码是否过期
	r, _ := cache.NewRedisCache(ctx).Get(utils.BuilderStr(service.UserEmail, ""))

	if r == nil {
		ctx.JSON(http.StatusOK, serializer.Err(serializer.CodeServerInternalError, "获取邮箱验证码失败", nil))
		return
	}

	rsp := service.AddUser()
	ctx.JSON(http.StatusOK, rsp)
}

// UserLogin 用户登录
func (uc *UserController) userLogin(ctx *gin.Context) {
	service := admin.NewUserLoginService()
	err := ctx.ShouldBindJSON(service)
	if err != nil {
		ctx.JSON(http.StatusOK, webutils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, service.Login())
	return
}

// UserInfo 用户信息
func (uc *UserController) userInfo(ctx *gin.Context) {
	service := admin.NewUserLoginService()
	ctx.JSON(http.StatusOK, service.UserInfo(ctx))
	return
}

// UserLogOut 用户登出
func (uc *UserController) userLogOut(ctx *gin.Context) {
	token := ctx.Request.Header.Get("X-Token")
	if token == "" {
		ctx.JSON(http.StatusOK, serializer.ParamErr("无效的请求", nil))
		return
	}
	service := admin.NewUserLoginService()
	ctx.JSON(http.StatusOK, service.UserLogOut(token))
}

func (uc *UserController) Build(rc *web.RouterCenter) {
	userGroup := rc.RG.Group(UserControllerGroup)
	userGroup.Handle("POST", "/login", uc.userLogin)
	userGroup.Handle("POST", "/register", uc.userRegister)
	userGroup.Handle("GET", "/user_code/:user_code", uc.checkUserExist)
	userGroup.Handle("GET", "/user_info", uc.userInfo)
	userGroup.Handle("POST", "/logout", uc.userLogOut)
}
