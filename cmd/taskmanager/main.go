package main

import (
	"taskmanager/internal/conf"
	"taskmanager/internal/taskManager"
	"taskmanager/internal/web"
	"taskmanager/internal/web/controller"
	"taskmanager/internal/web/middleware"
	"taskmanager/pkg/logger"
	"taskmanager/utils"
)

func init() {
	conf.LoadConf()
	logger.InitLogger()
	utils.InitDBSchema()
}

func main() {

	// 初始化任务
	go func() {
		tc := taskManager.GetDefaultTaskCenter()
		tc.StartWorker(10)
	}()

	// 初始化路由
	web.InitRouterCenter().
		Attach(
			middleware.NewCrossMiddleWare(),
			middleware.NewSessionMiddleWare(),
			middleware.NewErrorMiddleWare(),
			middleware.NewLoggerMiddleWare(),
		).
		Mount("api/", controller.NewMailController()).
		Mount("api/v1", controller.NewUserController()).
		Launch()

	utils.ServerNotify()
}
