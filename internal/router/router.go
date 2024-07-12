package router

import (
	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/middleware"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.Logger())

	// 设置 API 路由组
	api := r.Group("/api/v1")

	// 设置用户路由
	SetupUserRoutes(api)

	// 设置年级路由
	SetupGradeRoutes(api, db)
	SetupUnitRoutes(api, db)

	return r
}
