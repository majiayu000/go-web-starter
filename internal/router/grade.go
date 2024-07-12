package router

import (
	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/handlers"
	"github.com/majiayu000/gin-starter/internal/repositories"
	"github.com/majiayu000/gin-starter/internal/services"
	"gorm.io/gorm"
)

func SetupGradeRoutes(api *gin.RouterGroup, db *gorm.DB) {
	gradeRepo := repositories.NewGradeRepository(db)
	gradeService := services.NewGradeService(gradeRepo)
	gradeHandler := handlers.NewGradeHandler(gradeService)

	grades := api.Group("/grades")
	{
		grades.GET("/grades-by-language", gradeHandler.ListGradesByLanguage)
		grades.GET("", gradeHandler.ListGrades)
		grades.POST("", gradeHandler.CreateGrade)
		grades.GET("/:id", gradeHandler.GetGrade)
		grades.PUT("/:id", gradeHandler.UpdateGrade)
		grades.DELETE("/:id", gradeHandler.DeleteGrade)
	}
}
