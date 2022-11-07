package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"taskmanager/internal/conf"
	"time"
)

const (
	ManagerLog = iota // 服务日志
	GinLog            // gin框架日志
	TaskLog           // 任务日志
	AsynqLog          // 任务调度日志
)

var (
	managerLogger     *logrus.Logger
	ginLogger         *logrus.Logger
	taskLogger        *logrus.Logger
	asynqLogger       *logrus.Logger
	managerLumberjack *lumberjack.Logger
	ginLumberjack     *lumberjack.Logger
	taskLumberjack    *lumberjack.Logger
	asynqLumberjack   *lumberjack.Logger
)

func init() {
	managerLogger = logrus.New()
	ginLogger = logrus.New()
	taskLogger = logrus.New()
	asynqLogger = logrus.New()

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

	asynqLumberjack = &lumberjack.Logger{
		Filename:   conf.GetLogPath() + "asynq.log",
		MaxBackups: 5,
		MaxSize:    500,
		Compress:   true,
	}

}

func InitLogger() {
	AddHook(ManagerLog, &ManagerHook{
		file:     true,
		line:     true,
		function: true,
		levels:   logrus.AllLevels,
	})

	AddHook(TaskLog, &ManagerHook{
		file:     true,
		line:     true,
		function: true,
		levels:   logrus.AllLevels,
	})

	AddHook(GinLog, &ManagerHook{
		file:     false,
		line:     false,
		function: false,
		levels:   logrus.AllLevels,
	})

	AddHook(AsynqLog, &ManagerHook{
		file:     false,
		line:     false,
		function: false,
		levels:   logrus.AllLevels,
	})

	//设置日志格式化
	SetFormatter(ManagerLog, &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	SetFormatter(TaskLog, &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	SetFormatter(GinLog, &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	SetFormatter(AsynqLog, &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	level, err := logrus.ParseLevel(conf.GetLogLevel())
	if err != nil {
		panic(err)
	}
	//设置日志打印级别
	SetLevel(ManagerLog, level)
	SetLevel(TaskLog, level)
	SetLevel(GinLog, level)
	SetLevel(AsynqLog, level)

	SetOutput(ManagerLog, managerLumberjack)
	SetOutput(TaskLog, ginLumberjack)
	SetOutput(GinLog, taskLumberjack)
	SetOutput(AsynqLog, asynqLumberjack)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(logType int, formatter logrus.Formatter) {
	var std *logrus.Logger
	switch logType {
	case ManagerLog:
		std = managerLogger
	case TaskLog:
		std = taskLogger
	case GinLog:
		std = ginLogger
	case AsynqLog:
		std = asynqLogger
	default:
		panic("logger Type is not support")
	}
	std.Formatter = formatter
}

func SetOutput(logType int, out io.Writer) {
	var std *logrus.Logger
	switch logType {
	case ManagerLog:
		std = managerLogger
	case TaskLog:
		std = taskLogger
	case GinLog:
		std = ginLogger
	case AsynqLog:
		std = asynqLogger
	default:
		panic("logger Type is not support")
	}
	std.Out = io.MultiWriter(os.Stdout, out)
}

func SetLevel(logType int, level logrus.Level) {
	var std *logrus.Logger
	switch logType {
	case ManagerLog:
		std = managerLogger
	case TaskLog:
		std = taskLogger
	case GinLog:
		std = ginLogger
	case AsynqLog:
		std = asynqLogger
	default:
		panic("logger Type is not support")
	}
	std.SetLevel(level)
}

func Debug(format string, args ...interface{}) {
	managerLogger.Debugf(format, args...)
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

func GetAsynqLogger() *logrus.Logger {
	return asynqLogger
}
