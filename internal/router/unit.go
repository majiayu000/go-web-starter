package router

import (
	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/handlers"
	"github.com/majiayu000/gin-starter/internal/repositories"
	"github.com/majiayu000/gin-starter/internal/services"
	"gorm.io/gorm"
)

func SetupUnitRoutes(api *gin.RouterGroup, db *gorm.DB) {
	unitRepo := repositories.NewUnitRepository(db)
	unitService := services.NewUnitService(unitRepo)
	unitHandler := handlers.NewUnitHandler(unitService)

	units := api.Group("/units")
	{
		units.GET("", unitHandler.ListUnitsByLanguage)
	}
}
