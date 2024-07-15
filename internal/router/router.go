package router

import (
	"github.com/gin-gonic/gin"
	config "github.com/majiayu000/gin-starter/configs"
	"github.com/majiayu000/gin-starter/internal/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *redis.Client, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.Logger())
	logger := logrus.New()
	// 配置 logger
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	// 根据配置设置日志级别
	// if cfg.Debug {
	//     logger.SetLevel(logrus.DebugLevel)
	// } else {
	//     logger.SetLevel(logrus.InfoLevel)
	// }

	// 设置 API 路由组
	api := r.Group("/api/v1")

	// 设置用户路由
	SetupUserRoutes(api)

	// 设置年级路由
	SetupGradeRoutes(api, db)
	SetupUnitRoutes(api, db)
	SetupQuestionRoutes(api, db)
	SetupChatRoutes(api, db, redis, cfg, logger)

	return r
}
