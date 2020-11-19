package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ReqFromHTML ReqFromHtml
func ReqFromHTML(c *gin.Context) bool {
	req := c.Request
	reqAccept := strings.Split(req.Header.Get("Accept"), ",")[0]
	return reqAccept == "text/html"
}
