package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/conf"
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
	engine.Use(gin.Recovery()) // 默认中件件

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
