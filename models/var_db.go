package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库对象
var DB *gorm.DB

// Connect 链接数据库
func Connect() (err error) {
	db, err := gorm.Open(mysql.Open("root:sukeaiya@tcp(mysql:3306)/daily?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return
	}

	db.AutoMigrate(&ProjectModel{}, &ProjectLogModel{})

	DB = db
	return
}
