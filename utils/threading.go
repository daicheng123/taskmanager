package utils

import (
	"fmt"
	"taskmanager/pkg/logger"
)

func RunSafe(fn func(), errFunc func(err interface{})) {
	defer Recovery(errFunc)
	fn()
}

func Recovery(errFunc func(err interface{})) {
	if p := recover(); p != nil {
		errFunc(p)
	}
}

func RunSafeWithMsg(fn func(), errMsg string) {
	RunSafe(fn, func(err interface{}) {
		logger.Error(fmt.Sprintf("%s: %s", errMsg, err))
	})
}
