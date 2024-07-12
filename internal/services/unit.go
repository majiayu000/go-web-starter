package services

import (
	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/majiayu000/gin-starter/internal/repositories"
)

type UnitService struct {
	repo *repositories.UnitRepository
}

func (s *UnitService) ListUnitsByLanguage(lang string) (map[string][]models.Unit, error) {
	return s.repo.ListUnitsByLanguage(lang)
}

func NewUnitService(repo *repositories.UnitRepository) *UnitService {
	return &UnitService{repo: repo}
}
