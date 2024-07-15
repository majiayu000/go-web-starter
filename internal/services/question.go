// internal/services/question_service.go

package services

import (
	"context"

	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/majiayu000/gin-starter/internal/repositories"
)

type QuestionService struct {
	repo *repositories.QuestionRepository
}

func NewQuestionService(repo *repositories.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) GetQuestionsByGradeAndUnit(ctx context.Context, gradeID, unitNumber int, language string) ([]models.Question, error) {
	return s.repo.GetQuestionsByGradeAndUnit(ctx, gradeID, unitNumber, language)
}

func (s *QuestionService) GetQuestionByID(ctx context.Context, id int) (models.QuestionSentToModel, error) {
	return s.repo.GetQuestionByID(ctx, id)
}
