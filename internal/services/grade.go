package services

import (
	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/majiayu000/gin-starter/internal/repositories"
)

type GradeService struct {
	repo *repositories.GradeRepository
}

func NewGradeService(repo *repositories.GradeRepository) *GradeService {
	return &GradeService{
		repo: repo,
	}
}

func (s *GradeService) ListGradesByLanguage(lang string) ([]models.GradeResponse, error) {
	return s.repo.ListGradesByLanguage(lang)
}

func (s *GradeService) ListGrades() ([]models.Grade, error) {
	return s.repo.ListGrades()
}

func (s *GradeService) CreateGrade(grade *models.Grade) error {
	return s.repo.CreateGrade(grade)
}

func (s *GradeService) GetGradeByID(id int) (*models.Grade, error) {
	return s.repo.GetGradeByID(id)
}

func (s *GradeService) UpdateGrade(grade *models.Grade) error {
	return s.repo.UpdateGrade(grade)
}

func (s *GradeService) DeleteGrade(id int) error {
	return s.repo.DeleteGrade(id)
}
