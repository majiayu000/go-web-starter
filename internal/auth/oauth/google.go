package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/majiayu000/gin-starter/internal/types"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
	googleauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleProvider struct {
	config         *oauth2.Config
	sessionManager types.SessionManager
}

// Google login errors
var (
	ErrUnableToGetGoogleUser    = errors.New("google: unable to get Google User")
	ErrCannotValidateGoogleUser = errors.New("google: could not validate Google User")
)

func NewGoogleProvider(config map[string]string, sessionManager types.SessionManager) (*GoogleProvider, error) {
	clientID, ok := config["client_id"]
	if !ok {
		return nil, errors.New("Google client ID is missing")
	}

	clientSecret, ok := config["client_secret"]
	if !ok {
		return nil, errors.New("Google client secret is missing")
	}

	redirectURL, ok := config["redirect_url"]
	if !ok {
		return nil, errors.New("Google redirect URL is missing")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: googleOAuth2.Endpoint,
	}

	log.Printf("Google OAuth config initialized: %+v", oauthConfig)

	return &GoogleProvider{
		config:         oauthConfig,
		sessionManager: sessionManager,
	}, nil
}

func (p *GoogleProvider) GetAuthURL(state string) string {
	// 生成一个随机状态

	// 将状态保存到 session
	err := p.sessionManager.Set(context.Background(), "oauth_state:"+state, true, 10*time.Minute)
	if err != nil {
		log.Printf("Error saving state: %v", err)
		return ""
	}
	// 返回包含状态的 URL
	return p.config.AuthCodeURL(state)
}

func (p *GoogleProvider) Exchange(code string) (interface{}, error) {
	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (p *GoogleProvider) GetLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		state := generateRandomState()
		err := p.sessionManager.Set(c.Request.Context(), "oauth_state:"+state, true, 10*time.Minute)
		if err != nil {
			log.Printf("Error setting state in Redis: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize login"})
			return
		}

		url := p.GetAuthURL(state)
		log.Printf("Redirecting to auth URL: %s", url)
		c.Redirect(http.StatusFound, url)
	}
}

func generateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (p *GoogleProvider) GetCallbackHandler(successHandler http.Handler) gin.HandlerFunc {

	fmt.Println("CAaAAAAAAAAAA")
	return func(c *gin.Context) {
		// Implementation remains the same
	}
}

func (p *GoogleProvider) GetUserInfo(token interface{}) (*UserInfo, error) {
	oauthToken, ok := token.(*oauth2.Token)
	if !ok {
		return nil, errors.New("invalid token type for Google provider")
	}

	ctx := context.Background()
	client := p.config.Client(ctx, oauthToken)

	service, err := googleauth.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:    userInfo.Id,
		Name:  userInfo.Name,
		Email: userInfo.Email,
	}, nil
}
