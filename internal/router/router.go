// internal/router/router.go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/handlers"
	"github.com/majiayu000/gin-starter/internal/middleware"
)

// SetupRouter 初始化路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.Logger())
	// r.GET("/", handlers.HelloWorld)
	// 设置路由
	api := r.Group("/api")
	{
		api.GET("/user/:id", handlers.GetUser)
		// 在这里添加更多路由
	}

	// root := r.Group("/")
	// {
	// 	root.GET("/", handlers.HelloWorld)
	// }

	return r
}
