package repositories

import (
	"errors"

	"github.com/majiayu000/gin-starter/internal/models"
	"gorm.io/gorm"
)

type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{db: db}
}

func (r *UnitRepository) ListUnitsByLanguage(lang string) (map[string][]models.Unit, error) {
	var units []struct {
		models.Unit
		GradeCode string
	}
	var educationSystemID int

	switch lang {
	case "en":
		educationSystemID = 2
	case "zh":
		educationSystemID = 1
	default:
		return nil, errors.New("invalid language")
	}

	result := r.db.Table("t_math_demo_units").
		Joins("JOIN t_math_demo_textbooks ON t_math_demo_units.textbook_id = t_math_demo_textbooks.id").
		Joins("JOIN t_math_demo_grades ON t_math_demo_units.grade_id = t_math_demo_grades.id").
		Where("t_math_demo_textbooks.education_system_id = ?", educationSystemID).
		Select("t_math_demo_units.*, t_math_demo_grades.code as grade_code").
		Find(&units)

	if result.Error != nil {
		return nil, result.Error
	}

	gradeMap := make(map[string][]models.Unit)
	for _, unit := range units {
		gradeMap[unit.GradeCode] = append(gradeMap[unit.GradeCode], unit.Unit)
	}

	return gradeMap, nil
}
