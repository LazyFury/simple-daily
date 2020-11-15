package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinRecover Recover
func GinRecover(c *gin.Context) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			c.JSON(http.StatusNotFound, JSONError(err.Error(), err))
			return
		}
		if err, ok := r.(string); ok {
			c.JSON(http.StatusOK, JSONError(err, nil))
			return
		}
		if data, ok := r.(Result); ok {
			c.JSON(http.StatusOK, data)
			return
		}
		c.JSON(http.StatusInternalServerError, JSONError("error", nil))
	}

}
