package main

import (
	"taskmanager/internal/web"
	"taskmanager/internal/web/user"
)

func main() {
	web.InitRouterCenter().
		Attach(&web.CrossMiddleWare{}).
		Mount("v1", user.NewUserController()).
		Launch()
}
