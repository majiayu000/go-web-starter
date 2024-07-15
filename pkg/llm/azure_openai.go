package llm

import (
	"context"

	config "github.com/majiayu000/gin-starter/configs"
	"github.com/majiayu000/gin-starter/internal/interfaces"
	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/sashabaranov/go-openai"
)

type AzureOpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(cfg *config.Config) interfaces.LLMClient {
	apiKey := cfg.LLM.AzureOpenAI.APIKey
	endpoint := cfg.LLM.AzureOpenAI.Endpoint
	deploymentName := cfg.LLM.AzureOpenAI.DeploymentName
	// apiVersion := cfg.LLM.AzureOpenAI.APIVersion
	// config := openai.DefaultAzureConfig(authToken, url)
	config := openai.DefaultAzureConfig(apiKey, endpoint)
	if deploymentName != "" {
		config.AzureModelMapperFunc = func(model string) string {
			azureModelMapping := map[string]string{
				model: deploymentName,
			}
			return azureModelMapping[model]
		}
	}

	return &AzureOpenAIClient{
		client: openai.NewClientWithConfig(config),
	}
}

func (c *AzureOpenAIClient) CreateChatCompletion(ctx context.Context, messages []models.ChatMessage) (string, error) {
	// Convert models.ChatMessage to openai.ChatCompletionMessage
	openAIMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openAIMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: openAIMessages,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
