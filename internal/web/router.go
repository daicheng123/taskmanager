package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RouterCenter struct {
	*gin.Engine
	RG *gin.RouterGroup
}

func InitRouterCenter() *RouterCenter {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()

	// 心跳检测
	engine.Handle("GET", "/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, "OK")
	})

	return &RouterCenter{
		Engine: engine,
	}
}

// Attach 全局加载中间件
func (rc *RouterCenter) Attach(middleWares ...MiddleWare) *RouterCenter {
	rc.Use(func(context *gin.Context) {
		for _, middle := range middleWares {
			res := middle.OnRequest(context)
			if !res.Status {
				context.AbortWithStatusJSON(res.Code, res.Message)
			} else {
				context.Next()
			}
		}
	})
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
	rc.Run(fmt.Sprintf(":%d", 8080))
}
