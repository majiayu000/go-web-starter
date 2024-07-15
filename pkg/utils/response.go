// internal/utils/response.go

package utils

import (
	"github.com/gin-gonic/gin"
)

func ResponseJSON(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(code, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}
