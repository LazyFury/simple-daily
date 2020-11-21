package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinRecover Recover
func GinRecover(c *gin.Context) {

	if r := recover(); r != nil {
		result := JSON(http.StatusInternalServerError, "", nil)
		//普通错误
		if err, ok := r.(error); ok {
			log.Fatal(err)
			result.Message = err.Error()
			result.Data = err
		}
		//错误提示
		if err, ok := r.(string); ok {
			result.Message = err
		}
		//错误码
		if err, ok := r.(ErrCode); ok {
			result.Code = err
			result.Message = StatusText(err)
		} else if err, ok := r.(int); ok {
			result.Message = StatusText(ErrCode(err))
		}
		//完整错误类型
		if data, ok := r.(Result); ok {
			result = data
		}
		var code = http.StatusOK
		// if http.StatusText(int(result.Code)) != "" {
		// 	code = int(result.Code)
		// }
		//返回内容
		if ReqFromHTML(c) {
			c.HTML(code, "404.tmpl", result)
		} else {
			c.JSON(code, result)
		}

		log.Printf("\n\n\x1b[31m[Custom Debug Result]: URL:%s ;\nErr: %v \x1b[0m\n\n", c.Request.URL.RequestURI(), result)

		// panic("打断response继续写入内容")
		c.AbortWithStatus(http.StatusInternalServerError)

	}

}
