// internal/repositories/question_repository.go

package repositories

import (
	"context"

	"github.com/majiayu000/gin-starter/internal/models"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) GetQuestionsByGradeAndUnit(ctx context.Context, gradeID, unitNumber int, language string) ([]models.Question, error) {
	var questions []models.Question

	result := r.db.WithContext(ctx).
		Table("t_math_demo_questions").
		Where("grade = ? AND unit = ? AND language = ?", gradeID, unitNumber, language).
		Find(&questions)

	if result.Error != nil {
		return nil, result.Error
	}

	return questions, nil
}

func (r *QuestionRepository) GetQuestionByID(ctx context.Context, id int) (models.QuestionSentToModel, error) {
	var question models.QuestionSentToModel
	result := r.db.WithContext(ctx).
		Table("t_math_demo_questions").
		Select("id, content_plain as prompt_question, analysis_latex as solution").
		Where("id = ?", id).
		First(&question)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.QuestionSentToModel{}, nil // 没有找到问题
		}
		return models.QuestionSentToModel{}, result.Error
	}
	return question, nil
}
