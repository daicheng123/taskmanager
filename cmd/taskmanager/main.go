package main

import (
	"fmt"
	"taskmanager/utils"
)

//func init() {
//	conf.LoadConf()
//	logger.InitLogger()
//	utils.InitDBSchema()
//}
//
//func main() {
//
//	// 初始化任务
//	go func() {
//		tc := taskManager.GetDefaultTaskCenter()
//		tc.StartWorker(10)
//	}()
//
//	// 初始化路由
//	web.InitRouterCenter().
//		Attach(web.NewCrossMiddleWare()).
//		Mount("/api", common.NewCommonController()).
//		Mount("api/v1", user.NewUserController()).
//		Launch()
//
//	utils.ServerNotify()
//}
func main() {
	fmt.Println(utils.RandStringBytesMask(10))
}
