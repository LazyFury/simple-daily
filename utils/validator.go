package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// requiredParams 必填参数
var requiredParams validator.Func = func(fl validator.FieldLevel) bool {
	str, ok := fl.Field().Interface().(string)
	if ok {
		if str != "" {
			return true
		}
	}
	return false
}

// RegValidator 注册自定义验证器
func RegValidator() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = v.RegisterValidation("requiredParams", requiredParams)
	}
}
