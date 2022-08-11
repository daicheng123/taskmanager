package utils

import (
	"os"
	"os/signal"
	"syscall"
	"taskmanager/pkg/logger"
)

var ServerChan chan os.Signal

func init() {
	ServerChan = make(chan os.Signal)
}

func ShouldShutDown(err error) {
	if err != nil {
		logger.Error(err.Error())
		ServerChan <- os.Interrupt
	}
}

func ServerNotify() {
	signal.Notify(ServerChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ServerChan
	logger.Warning("执行退出操作")
}
