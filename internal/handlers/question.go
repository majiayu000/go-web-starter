// internal/handlers/question_handler.go

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/services"
	"github.com/majiayu000/gin-starter/pkg/utils"
)

type QuestionHandler struct {
	questionService *services.QuestionService
}

func NewQuestionHandler(questionService *services.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService}
}

var supportedLanguages = map[string]bool{
	"en": true,
	"zh": true,
	// 可以在这里添加更多支持的语言
}

func (h *QuestionHandler) GetQuestionsByGradeAndUnit(c *gin.Context) {
	gradeID, err := strconv.Atoi(c.Query("grade_id"))
	if err != nil {
		utils.ResponseJSON(c, http.StatusBadRequest, nil, "Invalid grade_id")
		return
	}

	unitNumber, err := strconv.Atoi(c.Query("unit_number"))
	if err != nil {
		utils.ResponseJSON(c, http.StatusBadRequest, nil, "Invalid unit_number")
		return
	}

	language := c.Query("language")
	if !supportedLanguages[language] {
		language = "zh" // 默认语言
	}

	questions, err := h.questionService.GetQuestionsByGradeAndUnit(c, gradeID, unitNumber, language)
	if err != nil {
		utils.ResponseJSON(c, http.StatusInternalServerError, nil, "Internal server error")
		return
	}

	utils.ResponseJSON(c, 0, questions, "Success")
}
