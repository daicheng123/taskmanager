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
	managerLogger     *logrus.Logger
	ginLogger         *logrus.Logger
	taskLogger        *logrus.Logger
	managerLumberjack *lumberjack.Logger
	ginLumberjack     *lumberjack.Logger
	taskLumberjack    *lumberjack.Logger
)

func init() {
	managerLogger = logrus.New()
	ginLogger = logrus.New()
	taskLogger = logrus.New()

	managerLumberjack = &lumberjack.Logger{
		Filename:   conf.GetLogPath() + "manager.log",
		MaxBackups: 5,
		MaxSize:    500,
		Compress:   true,
	}
	ginLumberjack = &lumberjack.Logger{
		Filename:   conf.GetLogPath() + "gin.log",
		MaxBackups: 5,
		MaxSize:    500,
		Compress:   true,
	}
	taskLumberjack = &lumberjack.Logger{
		Filename:   conf.GetLogPath() + "task.log",
		MaxBackups: 5,
		MaxSize:    500,
		Compress:   true,
	}
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
		file:     false,
		line:     false,
		function: false,
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

	SetOutput(consts.ManagerLog, managerLumberjack)
	SetOutput(consts.TaskLog, ginLumberjack)
	SetOutput(consts.GinLog, taskLumberjack)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(logType int, formatter logrus.Formatter) {
	var std *logrus.Logger
	switch logType {
	case consts.ManagerLog:
		std = managerLogger
	case consts.TaskLog:
		std = taskLogger
	case consts.GinLog:
		std = ginLogger
	default:
		panic("logger Type is not support")
	}
	std.Formatter = formatter
}

func SetOutput(logType int, out io.Writer) {
	var std *logrus.Logger
	switch logType {
	case consts.ManagerLog:
		std = managerLogger
	case consts.TaskLog:
		std = taskLogger
	case consts.GinLog:
		std = ginLogger
	default:
		panic("logger Type is not support")
	}
	std.Out = io.MultiWriter(os.Stdout, out)
}

func SetLevel(logType int, level logrus.Level) {
	var std *logrus.Logger
	switch logType {
	case consts.ManagerLog:
		std = managerLogger
	case consts.TaskLog:
		std = taskLogger
	case consts.GinLog:
		std = ginLogger
	default:
		panic("logger Type is not support")
	}
	std.SetLevel(level)
}

func Info(format string, args ...interface{}) {
	managerLogger.Infof(format, args...)
}

func Warning(format string, args ...interface{}) {
	managerLogger.Warningf(format, args...)
}

func Error(format string, args ...interface{}) {
	managerLogger.Errorf(format, args...)
}

func TaskInfo(format string, args ...interface{}) {
	managerLogger.Infof(format, args...)
}

func TaskWarning(format string, args ...interface{}) {
	managerLogger.Warningf(format, args...)
}

func TaskError(format string, args ...interface{}) {
	managerLogger.Errorf(format, args...)
}

func GinInfo(format string, args ...interface{}) {
	ginLogger.Infof(format, args...)
}

func GinERROR(format string, args ...interface{}) {
	ginLogger.Errorf(format, args...)
}
