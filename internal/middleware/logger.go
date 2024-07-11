// internal/middleware/logger.go
package middleware

import (
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.Logger()
}
