package database

import (
	"fmt"
	"strconv"

	config "github.com/majiayu000/gin-starter/configs"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.DB.Redis.Host, cfg.DB.Redis.Port)

	options := &redis.Options{
		Addr:     addr,
		Username: cfg.DB.Redis.User,
		Password: cfg.DB.Redis.Password,
		DB:       0, // 使用默认的数据库
	}

	if cfg.DB.Redis.Name != "" {
		// 如果指定了数据库名称，尝试将其转换为数据库索引
		dbIndex, err := strconv.Atoi(cfg.DB.Redis.Name)
		if err == nil {
			options.DB = dbIndex
		}
	}

	client := redis.NewClient(options)
	return client, nil

}
