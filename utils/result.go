package utils

import (
	"net/http"
	"time"
)

// Result Result
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	BuildBy time.Time   `json:"build_by"`
}

// BuildBy BuildBy
var BuildBy = time.Now()

// JSON JSON
func JSON(code int, message string, data interface{}) Result {
	return Result{
		Code:    code,
		Message: message,
		Data:    data,
		BuildBy: BuildBy,
	}
}

// JSONSuccess JSONSUCCESS
func JSONSuccess(message string, data interface{}) Result {
	if message == "" {
		message = "获取成功"
	}
	return JSON(http.StatusOK, message, data)
}

// JSONError JSONError
func JSONError(message string, data interface{}) Result {
	if message == "" {
		message = "error"
	}
	return JSON(http.StatusInternalServerError, message, data)
}
