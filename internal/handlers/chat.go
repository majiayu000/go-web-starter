package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	config "github.com/majiayu000/gin-starter/configs"
	"github.com/majiayu000/gin-starter/internal/interfaces"
	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/majiayu000/gin-starter/internal/services"
	"github.com/sirupsen/logrus"
)

type OpenAIHandler struct {
	config          *config.Config
	questionService *services.QuestionService
	chatService     *services.ChatService
	logger          *logrus.Logger
	redis           interfaces.ChatRepository
}

func NewOpenAIHandler(
	cfg *config.Config,
	questionService *services.QuestionService,
	chatService *services.ChatService,
	logger *logrus.Logger,
	redis interfaces.ChatRepository,
) *OpenAIHandler {
	return &OpenAIHandler{
		config:          cfg,
		questionService: questionService,
		chatService:     chatService,
		logger:          logger,
		redis:           redis,
	}
}

func (h *OpenAIHandler) CompletionHandler(c *gin.Context) {
	var request models.ChatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	response, err := h.chatService.ProcessChat(ctx, request)
	if err != nil {
		h.logger.WithError(err).Error("Failed to process chat")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process chat"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *OpenAIHandler) GetChatHistoryHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
		return
	}

	history, err := h.redis.GetChatHistory(c.Request.Context(), uuid)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get chat history")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *OpenAIHandler) ClearChatHistoryHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
		return
	}

	err := h.redis.ClearChatHistory(c.Request.Context(), uuid)
	if err != nil {
		h.logger.WithError(err).Error("Failed to clear chat history")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear chat history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat history cleared successfully"})
}

func (h *OpenAIHandler) AnswerHandler(c *gin.Context) {
	var request struct {
		QuestionID  int    `json:"questionId"`
		UUID        string `json:"uuid"`
		Language    string `json:language`
		UserMessage string `json:userMessage`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stopFlag := 0

	var newUUID string
	var messages []models.ChatMessage

	if request.UUID == "" {
		newUUID = uuid.New().String()
		question, err := h.questionService.GetQuestionByID(c.Request.Context(), request.QuestionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get question"})
			return
		}
		if (question == models.QuestionSentToModel{}) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
		messages, err = h.chatService.MakeUserMessage(question, request.Language)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user message"})
			return
		}
		// messages = append(h.config.PrimaryMessage, map[string]string{"role": "user", "content": userMessage})
	} else {
		newUUID = request.UUID
		// uuidKeyHistory = h.config.DB.Redis.Prefix + "::history::" + newUUID
		// uuidKeyFlag = h.config.DB.Redis.Prefix + "::flag::" + newUUID

		historyMessage, err := h.redis.GetChatHistory(c, newUUID)
		if err != nil {
			h.logger.WithError(err).Error("Failed to get history from Redis")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
			return
		}

		stopFlag, err := h.redis.GetFlag(c, newUUID)
		if err != nil {
			h.logger.WithError(err).Error("Failed to get stop flag from Redis")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stop flag"})
			return
		}

		fmt.Println("history missage is ", historyMessage)

		if stopFlag == 1 {
			c.JSON(http.StatusOK, gin.H{
				"uuid":     newUUID,
				"message":  "本题已经结束，请进入下一题",
				"stopFlag": stopFlag,
			})
			return
		}

		messages = append(historyMessage, models.ChatMessage{
			Role:    "user",
			Content: request.UserMessage})
		// fmt.Println("MESSAGE IS ", messages)
	}

	gptResponse, err := h.chatService.Chat(c, messages)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get GPT response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process chat"})
		return
	}

	// for _, keyword := range h.config.KeyWordList {
	// 	if strings.Contains(gptResponse, keyword) {
	// 		gptResponse = h.config.BadAnswer
	// 		break
	// 	}
	// }
	// clientIP := c.ClientIP()
	// h.logger.Infof("client host is %s, message is %+v", clientIP, messages)

	// h.logger.Infof("client host is %s, response is %s", clientIP, gptResponse)

	messages = append(messages, models.ChatMessage{
		Role:    "assistant",
		Content: gptResponse})

	// messagesJSON, _ := json.Marshal(messages)

	h.redis.SaveChatHistory(c, newUUID, messages)
	h.redis.SetFlag(c, newUUID, stopFlag)

	if strings.Contains(gptResponse, "总结") {
		stopFlag = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":     newUUID,
		"message":  gptResponse,
		"stopFlag": stopFlag,
	})
}
