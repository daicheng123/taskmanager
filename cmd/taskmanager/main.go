package main

import (
	"taskmanager/internal/conf"
	"taskmanager/internal/web"
	"taskmanager/internal/web/user"
)

func init() {
	conf.LoadConf()
}

func main() {
	web.InitRouterCenter().
		Attach(&web.CrossMiddleWare{}).
		Mount("v1", user.NewUserController()).
		Launch()
}
