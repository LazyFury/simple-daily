package models

// DB 全局数据库对象
var DB *GormDB = &GormDB{}

// Connect 链接数据库
func Connect(dsn string) (err error) {
	err = DB.ConnectMysql(dsn)
	if err != nil {
		return
	}

	DB.AutoMigrate(
		&ProjectModel{},
		&ProjectLogModel{},
		&UserModel{}, &FavoriteProjectModel{},
		&ProjectTodoModel{},
	)
	return
}
