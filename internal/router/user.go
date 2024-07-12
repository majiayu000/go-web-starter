package router

import (
	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/handlers"
)

func SetupUserRoutes(api *gin.RouterGroup) {
	api.GET("/user/:id", handlers.GetUser)
	// 在这里添加更多用户相关路由
}
