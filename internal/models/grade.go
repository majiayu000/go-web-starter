package models

type Grade struct {
	ID                int    `json:"id" gorm:"primaryKey"`
	EducationSystemID int    `json:"education_system_id"`
	Code              string `json:"code"`
	Name              string `json:"name"`
}

func (Grade) TableName() string {
	return "t_math_demo_grades" // 替换为您的实际表名
}

type GradeResponse struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type GradesWithUnits map[string][]Unit
