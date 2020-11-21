package routes

import (
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func setCsrfKey(c *gin.Context) (key string, err error) {
	key = utils.RandStringBytes(32)
	session := sessions.Default(c)
	session.Set("csrf", key)
	err = session.Save()
	return
}

func getCsrfKey(c *gin.Context) (csrf string, ok bool) {
	session := sessions.Default(c)
	csrf, ok = session.Get("csrf").(string)
	return
}
