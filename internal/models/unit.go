package models

type Unit struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	TextbookID  int    `json:"textbook_id"`
	UnitNumber  int    `json:"unit_number"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GradeID     int    `json:"grade_id" gorm:"column:grade_id"`
	GradeName   string `json:"grade_name" gorm:"column:grade_name"`
}

// TableName 指定模型对应的数据库表名
func (Unit) TableName() string {
	return "t_math_demo_units"
}

type UnitResponse struct {
	Code  string            `json:"code"`
	Units map[string][]Unit `json:"unit"`
}
