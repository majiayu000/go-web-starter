// internal/repositories/user.go
package repositories

import (
	"github.com/majiayu000/gin-starter/internal/models"
)

func GetUser(id string) (*models.User, error) {
	// 这里应该是数据库操作
	// 现在只是返回一个模拟的用户
	return &models.User{ID: id, Name: "John Doe"}, nil
}
