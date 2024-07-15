package services

import (
	"context"
	"fmt"

	config "github.com/majiayu000/gin-starter/configs"
	"github.com/majiayu000/gin-starter/internal/interfaces"
	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/majiayu000/gin-starter/pkg/llm"
	"github.com/majiayu000/gin-starter/pkg/utils"
)

type ChatService struct {
	repo           interfaces.ChatRepository
	azureOpenAICfg *config.Config
	azureOpenAI    interfaces.LLMClient // 假设您有一个 OpenAI 客户端接口
}

func NewChatService(repo interfaces.ChatRepository, cfg *config.Config) *ChatService {
	openAIClient := llm.NewOpenAIClient(cfg)

	return &ChatService{
		repo:           repo,
		azureOpenAICfg: cfg,
		azureOpenAI:    openAIClient,
	}
}

func (s *ChatService) ProcessChat(ctx context.Context, req models.ChatRequest) (*models.ChatResponse, error) {
	history, err := s.repo.GetChatHistory(ctx, req.UUID)
	if err != nil {
		return nil, err
	}

	messages := append(history, models.ChatMessage{Role: "user", Content: req.UserMessage})
	response, err := s.azureOpenAI.CreateChatCompletion(ctx, messages)
	if err != nil {
		return nil, err
	}

	messages = append(messages, models.ChatMessage{Role: "assistant", Content: response})
	if err := s.repo.SaveChatHistory(ctx, req.UUID, messages); err != nil {
		return nil, err
	}

	return &models.ChatResponse{
		UUID:     req.UUID,
		Message:  response,
		StopFlag: 0,
	}, nil
}

func (s *ChatService) Chat(ctx context.Context, messages []models.ChatMessage) (string, error) {
	// 将 models.ChatMessage 转换为 go-openai 包所需的格式
	// openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	// for i, msg := range messages {
	// 	openaiMessages[i] = openai.ChatCompletionMessage{
	// 		Role:    msg.Role,
	// 		Content: msg.Content,
	// 	}
	// }

	// 创建请求
	// request := openai.ChatCompletionRequest{
	// 	Model:    s.azureOpenAICfg.LLM.AzureOpenAI.DeploymentName, // 使用您配置的部署名称
	// 	Messages: openaiMessages,
	// }
	response, err := s.azureOpenAI.CreateChatCompletion(ctx, messages)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (s *ChatService) GetChatHistory(ctx context.Context, uuid string) ([]models.ChatMessage, error) {
	return s.repo.GetChatHistory(ctx, uuid)
}

func (s *ChatService) ClearChatHistory(ctx context.Context, uuid string) error {
	return s.repo.ClearChatHistory(ctx, uuid)
}

func (s *ChatService) GetSystemPrompt(lang string) string {
	switch lang {
	case "zh":
		return s.azureOpenAICfg.LLM.AzureOpenAI.SystemPromptZh
	case "en":
		return s.azureOpenAICfg.LLM.AzureOpenAI.SystemPromptEn
	default:
		return s.azureOpenAICfg.LLM.AzureOpenAI.SystemPromptDefault
	}
}

func (s *ChatService) MakeUserMessage(question models.QuestionSentToModel, language string) ([]models.ChatMessage, error) {
	var systemPrompt string
	// var primaryMessage []models.ChatMessage
	var userMessage models.ChatMessage

	fmt.Println("Language is:", language)

	if language == "en" {
		systemPrompt = s.azureOpenAICfg.LLM.AzureOpenAI.SystemPromptEn
		userMessage = models.ChatMessage{
			Role: "user",
			Content: fmt.Sprintf("Question:```%s```\nSolution process```%s```",
				question.PromptQuestion,
				question.Solution),
		}
	} else {
		systemPrompt = s.azureOpenAICfg.LLM.AzureOpenAI.SystemPromptZh
		userMessage = models.ChatMessage{
			Role: "user",
			Content: fmt.Sprintf("问题:```%s```\n解题过程:```%s```",
				question.PromptQuestion,
				question.Solution),
		}
	}

	messages := []models.ChatMessage{
		{Role: "system", Content: systemPrompt},
		userMessage,
	}

	utils.PrintMessages(messages)

	return messages, nil
}
