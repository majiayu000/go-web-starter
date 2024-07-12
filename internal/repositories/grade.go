package repositories

import (
	"errors"

	"github.com/majiayu000/gin-starter/internal/models"
	"gorm.io/gorm"
)

type GradeRepository struct {
	db *gorm.DB
}

func NewGradeRepository(db *gorm.DB) *GradeRepository {
	return &GradeRepository{
		db: db,
	}
}

func (r *GradeRepository) ListGradesByLanguage(lang string) ([]models.GradeResponse, error) {
	var grades []models.GradeResponse
	var educationSystemID int

	if lang == "en" {
		educationSystemID = 2
	} else if lang == "zh" {
		educationSystemID = 1
	} else {
		return nil, errors.New("invalid language")
	}

	result := r.db.Table("t_math_demo_grades").
		Select("id, code, name").
		Where("education_system_id = ?", educationSystemID).
		Find(&grades)

	return grades, result.Error
}

func (r *GradeRepository) ListGrades() ([]models.Grade, error) {
	var grades []models.Grade
	result := r.db.Find(&grades)
	return grades, result.Error
}

func (r *GradeRepository) CreateGrade(grade *models.Grade) error {
	return r.db.Create(grade).Error
}

func (r *GradeRepository) GetGradeByID(id int) (*models.Grade, error) {
	var grade models.Grade
	result := r.db.First(&grade, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &grade, nil
}

func (r *GradeRepository) UpdateGrade(grade *models.Grade) error {
	return r.db.Save(grade).Error
}

func (r *GradeRepository) DeleteGrade(id int) error {
	return r.db.Delete(&models.Grade{}, id).Error
}
