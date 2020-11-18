package models

// ModelType 数据库需要实现对应的便捷方法
type ModelType interface {
	Validator() (err error)
	Reset()
}
