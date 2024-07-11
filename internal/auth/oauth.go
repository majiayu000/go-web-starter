package auth

import (
	"errors"

	"github.com/majiayu000/gin-starter/internal/auth/oauth"
	"golang.org/x/oauth2"
)

type OAuthManager struct {
	providers map[string]oauth.Provider
}

func NewOAuthManager() *OAuthManager {
	return &OAuthManager{
		providers: make(map[string]oauth.Provider),
	}
}

func (m *OAuthManager) AddProvider(name string, provider oauth.Provider) {
	m.providers[name] = provider
}

func (m *OAuthManager) GetProvider(providerType string) (oauth.Provider, error) {
	provider, ok := m.providers[providerType]
	if !ok {
		return nil, errors.New("unknown provider type")
	}
	return provider, nil
}

func (m *OAuthManager) RemoveProvider(name string) {
	delete(m.providers, name)
}

func (m *OAuthManager) ListProviders() []string {
	providers := make([]string, 0, len(m.providers))
	for name := range m.providers {
		providers = append(providers, name)
	}
	return providers
}

func (m *OAuthManager) Exchange(providerType string, code string) (*oauth2.Token, error) {
	provider, ok := m.providers[providerType]
	if !ok {
		return nil, errors.New("unknown provider type")
	}

	tokenInterface, err := provider.Exchange(code)
	if err != nil {
		return nil, err
	}

	token, ok := tokenInterface.(*oauth2.Token)
	if !ok {
		return nil, errors.New("invalid token type")
	}

	return token, nil
}

func (m *OAuthManager) GetUserInfo(providerType string, token *oauth2.Token) (map[string]interface{}, error) {
	provider, ok := m.providers[providerType]
	if !ok {
		return nil, errors.New("unknown provider type")
	}

	userInfo, err := provider.GetUserInfo(token)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":    userInfo.ID,
		"name":  userInfo.Name,
		"email": userInfo.Email,
	}, nil
}
