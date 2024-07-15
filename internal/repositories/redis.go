package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/majiayu000/gin-starter/internal/interfaces"
	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/redis/go-redis/v9"
)

type redisChatRepository struct {
	client *redis.Client
	prefix string
}

func NewRedisChatRepository(client *redis.Client, prefix string) interfaces.ChatRepository {

	return &redisChatRepository{
		client: client,
		prefix: prefix,
	}
}

func (r *redisChatRepository) GetChatHistory(ctx context.Context, uuid string) ([]models.ChatMessage, error) {
	uuidKeyHistory := r.prefix + "::history::" + uuid
	fmt.Println("start get history")
	fmt.Println("key is ", uuid)
	data, err := r.client.Get(ctx, uuidKeyHistory).Bytes()
	if err != nil {
		if err == redis.Nil {
			return []models.ChatMessage{}, nil
		}
		return nil, err
	}

	// fmt.Println("history is ", data)

	var messages []models.ChatMessage
	err = json.Unmarshal(data, &messages)
	return messages, err
}

func (r *redisChatRepository) SaveChatHistory(ctx context.Context, uuid string, messages []models.ChatMessage) error {
	// fmt.Println("start save history", messages)
	key := r.prefix + "::history::" + uuid
	fmt.Println("redis key is", key)
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, 0).Err()
}

func (r *redisChatRepository) ClearChatHistory(ctx context.Context, uuid string) error {
	key := r.prefix + "::history::" + uuid
	return r.client.Del(ctx, key).Err()
}

func (r *redisChatRepository) GetFlag(ctx context.Context, uuid string) (int, error) {
	uuidKeyFlag := r.prefix + "::flag::" + uuid

	flagBytes, err := r.client.Get(ctx, uuidKeyFlag).Bytes()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	// 将 []byte 转换为 string，再转换为 int
	flagInt, err := strconv.Atoi(string(flagBytes))
	if err != nil {
		return 0, err
	}

	return flagInt, nil
}

func (r *redisChatRepository) SetFlag(ctx context.Context, uuid string, flag int) error {
	uuidKeyFlag := r.prefix + "::flag::" + uuid

	// 将 int 转换为 string
	flagStr := strconv.Itoa(flag)

	// 设置标志，这里假设我们想要设置一个过期时间（例如 24 小时）
	err := r.client.Set(ctx, uuidKeyFlag, flagStr, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}
