package main

import (
	"taskmanager/internal/conf"
	"taskmanager/internal/taskManager"
	"taskmanager/internal/web"
	"taskmanager/internal/web/common"
	"taskmanager/internal/web/user"
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
		Attach(web.NewCrossMiddleWare()).
		Mount("/api", common.NewCommonController()).
		Mount("api/v1", user.NewUserController()).
		Launch()

	utils.ServerNotify()
}
