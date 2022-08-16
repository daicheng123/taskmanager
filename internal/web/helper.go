package web

//Controller  控制器接口
type Controller interface {
	Build(rc *RouterCenter)
}
