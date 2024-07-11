package oauth

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/majiayu000/gin-starter/internal/types"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
	googleauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type AppleProvider struct {
	config         *oauth2.Config
	sessionManager types.SessionManager
}

func NewAppleProvider(config map[string]string, sessionManager types.SessionManager) (*AppleProvider, error) {
	clientID, ok := config["client_id"]
	if !ok {
		return nil, errors.New("Apple client ID is missing")
	}

	teamID, ok := config["team_id"]
	if !ok {
		return nil, errors.New("Apple client ID is missing")
	}

	keyID, ok := config["key_id"]
	if !ok {
		return nil, errors.New("Apple client secret is missing")
	}

	privateKey, ok := config["private_key"]
	if !ok {
		return nil, errors.New("Apple client ID is missing")
	}

	token, nil := generateToken(privateKey, teamID, clientID)

	redirectURL, ok := config["redirect_url"]
	if !ok {
		return nil, errors.New("Apple redirect URL is missing")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: token,
		KeyID:        keyID,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: googleOAuth2.Endpoint,
	}

	log.Printf("Apple OAuth config initialized: %+v", oauthConfig)

	return &AppleProvider{
		config:         oauthConfig,
		sessionManager: sessionManager,
	}, nil

	// return &Apple{
	// 	baseProvider: &baseProvider{
	// 		ctx:         context.Background(),
	// 		displayName: "Apple",
	// 		pkce:        true,
	// 		scopes:      []string{"name", "email"},
	// 		authUrl:     "https://appleid.apple.com/auth/authorize",
	// 		tokenUrl:    "https://appleid.apple.com/auth/token",
	// 	},
	// 	jwksUrl: "https://appleid.apple.com/auth/keys",
	// }
}

func (p *AppleProvider) GetAuthURL(state string) string {
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

func generateToken(file string, clientID string, teamID string) (string, error) {
	privateKeyData, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("无法读取私钥文件: %v", err)
	}

	// 解析私钥
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		log.Fatalf("无法解码 PEM 块")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("无法解析私钥: %v", err)
	}

	ecdsaKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		log.Fatalf("私钥不是有效的 ECDSA 私钥")
	}

	// 设置 JWT 声明
	claims := jwt.MapClaims{
		"iss": teamID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 180).Unix(), // 6个月后过期
		"aud": "https://appleid.apple.com",
		"sub": clientID,
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = "YOUR_KEY_ID"

	// 签名并获取完整的编码后的字符串 token
	tokenString, err := token.SignedString(ecdsaKey)
	if err != nil {
		log.Fatalf("无法签名 token: %v", err)
	}

	return tokenString, nil

}

func (p *AppleProvider) Exchange(code string) (interface{}, error) {
	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (p *AppleProvider) GetLoginHandler() gin.HandlerFunc {
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

func (p *AppleProvider) GetCallbackHandler(successHandler http.Handler) gin.HandlerFunc {

	fmt.Println("CAaAAAAAAAAAA")
	return func(c *gin.Context) {
		// Implementation remains the same
	}
}

func (p *AppleProvider) GetUserInfo(token interface{}) (*UserInfo, error) {
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
