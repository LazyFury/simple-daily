package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库对象
var DB *gorm.DB

// Connect 链接数据库
func Connect(dsn string) (err error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return
	}

	db.AutoMigrate(
		&ProjectModel{},
		&ProjectLogModel{},
		&UserModel{}, &FavoriteProjectModel{},
		&ProjectTodoModel{},
	)

	DB = db

	return
}
