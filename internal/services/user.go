// internal/services/user.go
package services

import (
	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/majiayu000/gin-starter/internal/repositories"
)

func GetUser(id string) (*models.User, error) {
	return repositories.GetUser(id)
}
