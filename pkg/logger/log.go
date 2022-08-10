package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"taskmanager/internal/conf"
	"taskmanager/internal/consts"
	"time"
)

var (
	ManagerLogger     *logrus.Logger
	GinLogger         *logrus.Logger
	TaskLogger        *logrus.Logger
	ManagerLumberjack *lumberjack.Logger
	GinLumberjack     *lumberjack.Logger
	TaskLumberjack    *lumberjack.Logger
)

func init() {
	ManagerLogger = new(logrus.Logger)
	GinLogger = new(logrus.Logger)
	TaskLogger = new(logrus.Logger)

	ManagerLumberjack = new(lumberjack.Logger)
	GinLumberjack = new(lumberjack.Logger)
	TaskLumberjack = new(lumberjack.Logger)
}

func InitLogger() {
	AddHook(consts.ManagerLog, &ManagerHook{
		file:     true,
		line:     true,
		function: true,
		levels:   logrus.AllLevels,
	})

	AddHook(consts.TaskLog, &ManagerHook{
		file:     true,
		line:     true,
		function: true,
		levels:   logrus.AllLevels,
	})

	AddHook(consts.GinLog, &ManagerHook{
		file:     true,
		line:     true,
		function: true,
		levels:   logrus.AllLevels,
	})

	//设置日志格式化
	SetFormatter(consts.ManagerLog, &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	SetFormatter(consts.TaskLog, &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	SetFormatter(consts.GinLog, &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	level, err := logrus.ParseLevel(conf.GetLogLevel())
	if err != nil {
		panic(err)
	}
	//设置日志打印级别
	SetLevel(consts.ManagerLog, level)
	SetLevel(consts.TaskLog, level)
	SetLevel(consts.GinLog, level)

	SetOutput(consts.ManagerLog, ManagerLumberjack)
	SetOutput(consts.TaskLog, GinLumberjack)
	SetOutput(consts.GinLog, TaskLumberjack)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(logType int, formatter logrus.Formatter) {
	var std *logrus.Logger
	switch logType {
	case consts.ManagerLog:
		std = ManagerLogger
	case consts.TaskLog:
		std = TaskLogger
	case consts.GinLog:
		std = GinLogger
	default:
		panic("logger Type is not support")
	}
	std.Formatter = formatter
}

func SetOutput(logType int, out io.Writer) {
	var std *logrus.Logger
	switch logType {
	case consts.ManagerLog:
		std = ManagerLogger
	case consts.TaskLog:
		std = TaskLogger
	case consts.GinLog:
		std = GinLogger
	default:
		panic("logger Type is not support")
	}
	std.Out = io.MultiWriter(os.Stdout, out)
}

func SetLevel(logType int, level logrus.Level) {
	var std *logrus.Logger
	switch logType {
	case consts.ManagerLog:
		std = ManagerLogger
	case consts.TaskLog:
		std = TaskLogger
	case consts.GinLog:
		std = GinLogger
	default:
		panic("logger Type is not support")
	}
	std.SetLevel(level)
}

func Info(format string, args ...interface{}) {
	ManagerLogger.Infof(format, args...)
}

func Warning(format string, args ...interface{}) {
	ManagerLogger.Warningf(format, args...)
}

func Error(format string, args ...interface{}) {
	ManagerLogger.Errorf(format, args...)
}
