package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"taskmanager/internal/conf"
	modelscommon "taskmanager/internal/models/common"
	"taskmanager/internal/web/middleware"
	"taskmanager/pkg/logger"
)

type RouterCenter struct {
	*gin.Engine
	RG *gin.RouterGroup
}

func InitRouterCenter() *RouterCenter {
	gin.SetMode(conf.GetWebMode())
	engine := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册 models.CustomTime 类型的自定义校验规则
		v.RegisterCustomTypeFunc(modelscommon.ValidateJSONDateType, modelscommon.CustomTime{})
	}

	// 心跳检测
	engine.Handle("GET", "/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, "OK")
	})

	return &RouterCenter{
		Engine: engine,
	}
}

// Attach 全局加载中间件
func (rc *RouterCenter) Attach(middleWares ...middleware.MiddleWare) *RouterCenter {
	rc.Engine.Use(gin.Recovery())
	for _, m := range middleWares {
		rc.Engine.Use(m.OnRequest())
	}
	return rc
}

//Mount  挂载控制器
func (rc *RouterCenter) Mount(group string, controllers ...Controller) *RouterCenter {
	rc.RG = rc.Group(group)
	for _, controller := range controllers {
		controller.Build(rc)
	}
	return rc
}

// Launch 启动
func (rc *RouterCenter) Launch() { // 需要把Index 和 UserIndex 两个控制器传递进来
	port := conf.GetWebPort()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: rc.Engine}

	go func() {
		logger.Info("init server on %d", port)
		if err := srv.ListenAndServe(); err != nil {
			return
		}
	}()
}
