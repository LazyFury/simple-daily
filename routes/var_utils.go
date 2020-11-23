package routes

import (
	"log"
	"net/http"

	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func setCsrfKey(c *gin.Context) (key string, err error) {
	key = utils.RandStringBytes(32)
	session := sessions.Default(c)
	session.Set("csrf", key)
	err = session.Save()
	return
}

func getCsrfKey(c *gin.Context) (csrf string, ok bool, err error) {
	session := sessions.Default(c)
	csrf, ok = session.Get("csrf").(string)
	session.Delete("csrf")
	err = session.Save()
	return
}

func csrf(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		_csrf, err := setCsrfKey(c)
		if err != nil {
			panic(err)
		}
		c.Set("csrf", _csrf)
	}
	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete {
		form := struct {
			Csrf string `form:"csrf" binding:"required"`
		}{}

		if c.ContentType() == binding.MIMEJSON {
			log.Print("csrf binding")
			if err := utils.ShouldBindBodyWith(c, &form, binding.JSON); err != nil {
				panic(utils.JSONError("验证失败 请刷新页面重试", err))
			}
		} else {
			b := binding.Default(c.Request.Method, c.ContentType())
			if err := c.ShouldBindWith(&form, b); err != nil {
				panic(utils.JSONError("验证失败 请刷新页面重试", err))
			}
		}

		csrf, ok, err := getCsrfKey(c)
		if !ok || csrf != form.Csrf || err != nil {
			panic(utils.JSONError("验证失败", err))
		}
	}

	c.Next()
}
