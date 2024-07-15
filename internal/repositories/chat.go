package repositories

import (
	"github.com/majiayu000/gin-starter/internal/interfaces"
	"github.com/redis/go-redis/v9"
)

// NewChatRepository 创建一个新的 ChatRepository 实例
func NewChatRepository(redisURL string) interfaces.ChatRepository {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	return &redisChatRepository{client: client}
}
