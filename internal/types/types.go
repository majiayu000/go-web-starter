// internal/types/types.go
package types

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type SessionManager interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	GetSession(c *gin.Context) (*oauth2.Token, map[string]interface{}, error)
}
