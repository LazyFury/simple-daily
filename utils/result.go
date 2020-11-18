package utils

import (
	"time"
)

// Result Result
type Result struct {
	Code    ErrCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	BuildBy time.Time   `json:"build_by"`
}

// ErrCode 错误码类型
type ErrCode int

const (
	// Success Success
	Success ErrCode = 1
	// Errors 失败
	Errors ErrCode = -1
)
const (
	// LoginSuccess 登陆成功
	LoginSuccess ErrCode = iota + 100
)
const (
	// AuthenError 认证失败
	AuthenError ErrCode = -iota - 100
	// NotFound 没有数据
	NotFound
)

// ErrorCodeText 错误提示
var ErrorCodeText = map[ErrCode]string{
	Success:      "获取成功",
	Errors:       "遇到错误",
	LoginSuccess: "登陆成功",
	AuthenError:  "登陆超时",
	NotFound:     "没有数据",
}

// BuildBy BuildBy
var BuildBy = time.Now()

// StatusText StatusText
func StatusText(code ErrCode) string {
	msg := ErrorCodeText[code]
	if msg == "" {
		msg = "未知错误码"
	}
	return msg
}

// JSON JSON
func JSON(code ErrCode, message string, data interface{}) Result {
	if message == "" {
		message = StatusText(code)
	}
	return Result{
		Code:    code,
		Message: message,
		Data:    data,
		BuildBy: BuildBy,
	}
}

// JSONSuccess JSONSUCCESS
func JSONSuccess(message string, data interface{}) Result {
	return JSON(Success, message, data)
}

// JSONError JSONError
func JSONError(message string, data interface{}) Result {
	return JSON(Errors, message, data)
}
