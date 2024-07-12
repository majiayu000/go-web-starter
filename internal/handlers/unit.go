package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/services"
)

type UnitHandler struct {
	service *services.UnitService
}

func NewUnitHandler(service *services.UnitService) *UnitHandler {
	return &UnitHandler{service: service}
}

func (h *UnitHandler) ListUnitsByLanguage(c *gin.Context) {
	lang := c.Query("lang")
	if lang != "en" && lang != "zh" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language parameter. Use 'en' or 'zh'."})
		return
	}

	units, err := h.service.ListUnitsByLanguage(lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, units)
}
