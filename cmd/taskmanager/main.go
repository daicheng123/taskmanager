package main

import (
	"taskmanager/internal/conf"
	modelutils "taskmanager/internal/models/utils"
	"taskmanager/internal/web"
	"taskmanager/internal/web/controller"
	"taskmanager/internal/web/middleware"
	"taskmanager/pkg/logger"
	"taskmanager/utils"
)

func init() {
	conf.LoadConf()
	logger.InitLogger()
	modelutils.InitDBSchema()
}

func main() {

	// 初始化任务
	//go func() {
	//	tc := taskManager.GetDefaultTaskCenter()
	//	tc.StartWorker(10)
	//}()

	// 初始化路由
	web.InitRouterCenter().
		Attach(
			middleware.NewCrossMiddleWare(),
			middleware.NewSessionMiddleWare(),
			middleware.NewErrorMiddleWare(),
			middleware.NewLoggerMiddleWare(),
		).
		// 邮件相关接口
		Mount("api/", controller.NewMailController()).
		// 用户相关接口
		Mount("api/v1", controller.NewUserController()).
		// 标签相关接口
		Mount("api/v1", controller.NewTagsController()).
		// 执行器相关接口
		Mount("api/v1", controller.NewExecutorController()).
		// 操作脚本相关接口
		Mount("api/v1", controller.NewScriptController()).
		// 危险命令相关接口
		Mount("api/v1", controller.NewDangerCmdController()).
		Launch()
	utils.ServerNotify()

}
