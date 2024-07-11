package oauth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/types"
)

type UserInfo struct {
	ID    string
	Name  string
	Email string
}

type Provider interface {
	GetAuthURL(state string) string
	Exchange(code string) (interface{}, error)
	GetLoginHandler() gin.HandlerFunc
	GetCallbackHandler(successHandler http.Handler) gin.HandlerFunc
	GetUserInfo(token interface{}) (*UserInfo, error)
}

type ProviderType string

const (
	Google   ProviderType = "google"
	Apple    ProviderType = "apple"
	Facebook ProviderType = "facebook"
)

func NewProvider(providerType ProviderType, config map[string]string, sessionManager types.SessionManager) (Provider, error) {
	switch providerType {
	case Google:
		return NewGoogleProvider(config, sessionManager)
	case Apple:
		return NewAppleProvider(config, sessionManager)

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
