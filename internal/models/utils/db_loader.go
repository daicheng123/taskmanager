package utils

import (
	"os"
	"taskmanager/internal/models"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/store"
)

func InitDBSchema() {
	autoMigrate()
}

//autoMigrate 自动根据数据模型建表，表名为实体名的蛇形表示
func autoMigrate() {
	db, err := store.GetOrCreateDBOperator()
	if err != nil {
		logger.Error("获取数据库连接失败")
		os.Exit(1)
	}

	err = db.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(
			&models.UserModel{},
			&models.SessionModel{},
			&models.Tag{},
			&models.Executor{},
			&models.ExecutorAccount{},
			&models.Script{},
			&models.DangerousCmd{},
			&models.ScriptAudit{},
		)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
