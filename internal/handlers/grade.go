package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/majiayu000/gin-starter/internal/services"
)

type GradeHandler struct {
	gradeService *services.GradeService
}

func NewGradeHandler(gradeService *services.GradeService) *GradeHandler {
	return &GradeHandler{
		gradeService: gradeService,
	}
}

func (h *GradeHandler) ListGradesByLanguage(c *gin.Context) {
	lang := c.Query("lang")
	if lang != "en" && lang != "zh" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language parameter. Use 'en' or 'zh'."})
		return
	}

	grades, err := h.gradeService.ListGradesByLanguage(lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}

func (h *GradeHandler) ListGrades(c *gin.Context) {
	grades, err := h.gradeService.ListGrades()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grades)
}

func (h *GradeHandler) CreateGrade(c *gin.Context) {
	var grade models.Grade
	if err := c.ShouldBindJSON(&grade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.gradeService.CreateGrade(&grade); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, grade)
}

func (h *GradeHandler) GetGrade(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade ID"})
		return
	}

	grade, err := h.gradeService.GetGradeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grade not found"})
		return
	}
	c.JSON(http.StatusOK, grade)
}

func (h *GradeHandler) UpdateGrade(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade ID"})
		return
	}

	var grade models.Grade
	if err := c.ShouldBindJSON(&grade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grade.ID = id

	if err := h.gradeService.UpdateGrade(&grade); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grade)
}

func (h *GradeHandler) DeleteGrade(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade ID"})
		return
	}

	if err := h.gradeService.DeleteGrade(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Grade deleted successfully"})
}
