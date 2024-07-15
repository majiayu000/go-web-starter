package models

import "github.com/gofrs/uuid"

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	UUID        string `json:"uuid"`
	UserMessage string `json:"user_message"`
}

type ChatResponse struct {
	UUID     string `json:"uuid"`
	Message  string `json:"message"`
	StopFlag int    `json:"stop_flag"`
}

type UserMessage struct {
	QuestionID  *int                   `json:"questionId,omitempty"`
	UserMessage *string                `json:"userMessage,omitempty"`
	UserDict    map[string]interface{} `json:"userDict,omitempty"`
	UUID        *uuid.UUID             `json:"uuid,omitempty"`
	StopFlag    *int                   `json:"stopFlag,omitempty"`
}
