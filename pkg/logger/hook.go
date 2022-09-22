package logger

import (
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
	"sync"
)

const (
	skipDepth = 3
)

type ManagerHook struct {
	mu       sync.Mutex
	file     bool // 是否打印文件名称
	line     bool // 是否打印行号
	function bool // 是否打印函数
	levels   []logrus.Level
}

func (ths *ManagerHook) Levels() []logrus.Level {
	return ths.levels
}
func (ths *ManagerHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 30)
	n := runtime.Callers(skipDepth, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		funName := frame.Function
		if !strings.Contains(funName, "sirupsen/logrus") &&
			!strings.Contains(funName, "taskmanager/pkg/logger") &&
			!strings.Contains(funName, "runtime") {
			//fmt.Println(funName)
			if ths.file {
				ths.mu.Lock()
				entry.Data["file"] = path.Base(frame.File)
				ths.mu.Unlock()
			}
			if ths.line {
				ths.mu.Lock()
				entry.Data["line"] = frame.Line
				ths.mu.Unlock()
			}
			if ths.function {
				ths.mu.Lock()
				entry.Data["function"] = frame.Function
				ths.mu.Unlock()
			}
			break // 找到第一个与日志无关的函数指针后退出
		}
		if !more {
			break
		}
	}
	return nil
}

func AddHook(logType int, hook logrus.Hook) {
	var std *logrus.Logger
	switch logType {
	case ManagerLog:
		std = managerLogger
	case GinLog:
		std = ginLogger
	case TaskLog:
		std = taskLogger
	default:
		panic("logger Type is not support")
	}
	std.Hooks.Add(hook)
}
