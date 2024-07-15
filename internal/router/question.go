// internal/router/question_router.go

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/handlers"
	"github.com/majiayu000/gin-starter/internal/repositories"
	"github.com/majiayu000/gin-starter/internal/services"
	"gorm.io/gorm"
)

func SetupQuestionRoutes(api *gin.RouterGroup, db *gorm.DB) {
	questionRepo := repositories.NewQuestionRepository(db)
	questionService := services.NewQuestionService(questionRepo)
	questionHandler := handlers.NewQuestionHandler(questionService)

	questionGroup := api.Group("/questions")
	{
		questionGroup.GET("", questionHandler.GetQuestionsByGradeAndUnit)
	}
}
