package store

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"sync"
	"taskmanager/internal/conf"
	"taskmanager/pkg/logger"
)

var (
	dbOperator *gorm.DB
	dbOnce     sync.Once
)

func newDBOperator() (err error) {
	dbOnce.Do(func() {
		dsn := fmt.Sprintf(
			`%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local`,
			conf.GetDBUser(),
			conf.GetDbPassword(),
			conf.GetDBAddress(),
			conf.GetDBPort(),
			conf.GetDbName())
		fmt.Printf("dsn: %s", dsn)
		logLevel := gormLogger.Warn

		if conf.IsDebugMode() {
			logLevel = gormLogger.Info
		}

		dbOperator, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(logLevel),
		})

		if err != nil {
			logger.Error(err.Error())
			return
		}
		sqlDB, err := dbOperator.DB()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		sqlDB.SetConnMaxIdleTime(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(20)
	})
	return
}

// GetOrCreateDBOperator 获取数据库会话
func GetOrCreateDBOperator() (*gorm.DB, error) {
	if dbOperator == nil {
		if err := newDBOperator(); err != nil {
			return nil, fmt.Errorf("创建数据库连接失败: %s", err)
		}
	}
	return dbOperator, nil
}

//Execute  执行数据访问操作
func Execute(dbOperate func(*gorm.DB) *gorm.DB) (*gorm.DB, error) {
	dbSession, err := GetOrCreateDBOperator()
	if err != nil {
		return nil, err
	}
	if dbSession == nil {
		return dbSession, fmt.Errorf("创建数据库连接失败")
	}
	result := dbOperate(dbSession)
	return result, result.Error
}
