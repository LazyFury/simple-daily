package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库对象
var DB *gorm.DB

// Connect 链接数据库
func Connect() (err error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return
	}

	db.AutoMigrate(&ProjectModel{}, &ProjectLogModel{})

	DB = db
	return
}
