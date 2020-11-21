package models

// ModelType 数据库需要实现对应的便捷方法
type ModelType interface {
	//Validator 验证模型
	//常用的非空验证 友好提示
	Validator() (err error)

	//TableName
	//邀请数据模型必须提供tablename
	//gorm依次自动迁移建表
	//查询数据比如join操作可以便捷拼接表名
	TableName() string
}
