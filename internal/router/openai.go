package router

import (
	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/handlers"
	"github.com/majiayu000/gin-starter/internal/repositories"
	"github.com/majiayu000/gin-starter/internal/services"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	config "github.com/majiayu000/gin-starter/configs"

	"github.com/sirupsen/logrus"
)

func SetupChatRoutes(api *gin.RouterGroup, db *gorm.DB, redis *redis.Client, cfg *config.Config, logger *logrus.Logger) {
	chatRepo := repositories.NewRedisChatRepository(redis, cfg.DB.Redis.Prefix)
	questionRepo := repositories.NewQuestionRepository(db)
	questionService := services.NewQuestionService(questionRepo)

	chatService := services.NewChatService(chatRepo, cfg)

	chatHandler := handlers.NewOpenAIHandler(
		cfg,
		questionService,
		chatService,
		logger,
		chatRepo,
	)

	chatGroup := api.Group("/chat")
	{
		chatGroup.POST("", chatHandler.AnswerHandler)
		chatGroup.GET("/:uuid", chatHandler.GetChatHistoryHandler)
		chatGroup.DELETE("/:uuid", chatHandler.ClearChatHistoryHandler)
	}
}
