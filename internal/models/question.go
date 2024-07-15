package models

import "time"

type Question struct {
	ID              int       `json:"id"`
	SubjectID       int       `json:"subject_id"`
	Grade           int       `json:"grade"`
	Unit            int       `json:"unit"`
	QuestionType    int       `json:"question_type"`
	Difficulty      int       `json:"difficulty"`
	Content         string    `json:"content"`
	ContentPlain    string    `json:"content_plain"`
	ContentLatex    string    `json:"content_latex"`
	Option          string    `json:"option"`
	Answer          string    `json:"answer"`
	AnswerLatex     string    `json:"answer_latex"`
	Analysis        string    `json:"analysis"`
	AnalysisLatex   string    `json:"analysis_latex"`
	KnowledgePoints string    `json:"knowledge_points"`
	Source          string    `json:"source"`
	SourceImage     string    `json:"source_image"`
	CreatedTime     time.Time `json:"created_time"`
	UpdatedTime     time.Time `json:"updated_time"`
	Status          string    `json:"status"`
	Language        string    `json:"language"`
}

func (Question) TableName() string {
	return "t_math_demo_questions"
}

type QuestionSentToModel struct {
	ID             int    `gorm:"column:id"`
	PromptQuestion string `gorm:"column:prompt_question"`
	Solution       string `gorm:"column:solution"`
}
