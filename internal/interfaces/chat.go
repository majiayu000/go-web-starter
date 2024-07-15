package interfaces

import (
	"context"

	"github.com/majiayu000/gin-starter/internal/models"
)

type ChatRepository interface {
	GetChatHistory(ctx context.Context, uuid string) ([]models.ChatMessage, error)
	SaveChatHistory(ctx context.Context, uuid string, messages []models.ChatMessage) error
	ClearChatHistory(ctx context.Context, uuid string) error
	GetFlag(ctx context.Context, uuid string) (int, error)
	SetFlag(ctx context.Context, uuid string, flag int) error
}

type ChatService interface {
	Chat(ctx context.Context, messages []models.ChatMessage) (string, error)
	ProcessChat(ctx context.Context, req models.ChatRequest) (*models.ChatResponse, error)
	GetChatHistory(ctx context.Context, uuid string) ([]models.ChatMessage, error)
	ClearChatHistory(ctx context.Context, uuid string) error
	GetQuestion(questionID int) (*models.QuestionSentToModel, error)
	MakeUserMessage(question models.QuestionSentToModel, language string) ([]models.ChatMessage, error)
}

type LLMClient interface {
	CreateChatCompletion(ctx context.Context, messages []models.ChatMessage) (string, error)
}

type OpenAIClient interface {
	CreateChatCompletion(ctx context.Context, messages []models.ChatMessage) (string, error)
}
