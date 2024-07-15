package services

import (
	"context"

	"github.com/majiayu000/gin-starter/internal/interfaces"
	"github.com/majiayu000/gin-starter/internal/models"
)

type RedisService struct {
	repo interfaces.ChatRepository
}

func NewRedisService(repo interfaces.ChatRepository) *RedisService {
	return &RedisService{
		repo: repo,
	}
}

// GetChatHistory 获取聊天历史
func (s *RedisService) GetChatHistory(ctx context.Context, uuid string) ([]models.ChatMessage, error) {
	return s.repo.GetChatHistory(ctx, uuid)
}

// SaveChatHistory 保存聊天历史
func (s *RedisService) SaveChatHistory(ctx context.Context, uuid string, messages []models.ChatMessage) error {
	return s.repo.SaveChatHistory(ctx, uuid, messages)
}

// ClearChatHistory 清除聊天历史
func (s *RedisService) ClearChatHistory(ctx context.Context, uuid string) error {
	return s.repo.ClearChatHistory(ctx, uuid)
}

// AddMessageToHistory 添加新消息到历史记录
func (s *RedisService) AddMessageToHistory(ctx context.Context, uuid string, message models.ChatMessage) error {
	history, err := s.GetChatHistory(ctx, uuid)
	if err != nil {
		return err
	}

	history = append(history, message)
	return s.SaveChatHistory(ctx, uuid, history)
}

// GetLastNMessages 获取最后N条消息
func (s *RedisService) GetLastNMessages(ctx context.Context, uuid string, n int) ([]models.ChatMessage, error) {
	history, err := s.GetChatHistory(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if len(history) <= n {
		return history, nil
	}

	return history[len(history)-n:], nil
}

// TruncateHistory 截断历史记录至指定长度
func (s *RedisService) TruncateHistory(ctx context.Context, uuid string, maxLength int) error {
	history, err := s.GetChatHistory(ctx, uuid)
	if err != nil {
		return err
	}

	if len(history) > maxLength {
		history = history[len(history)-maxLength:]
		return s.SaveChatHistory(ctx, uuid, history)
	}

	return nil
}
