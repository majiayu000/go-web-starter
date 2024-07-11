package auth

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/majiayu000/gin-starter/internal/types"
	"github.com/redis/go-redis/v9"

	"golang.org/x/oauth2"
)

type SessionInfo struct {
	AccessToken string
	TokenType   string
	Expiry      time.Time
	UserID      string
	UserName    string
	Email       string
	// 其他必要的用户信息字段
}

type Session struct {
	AccessToken string
	TokenType   string
	Expiry      time.Time
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type SessionManager struct {
	redisClient *redis.Client
}

func NewSessionManager(redisAddr, redisPassword string, redisDB int) types.SessionManager {
	return &SessionManager{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		}),
	}
}

func (sm *SessionManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return sm.redisClient.Set(ctx, key, jsonValue, expiration).Err()
}

func (sm *SessionManager) Get(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := sm.redisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonValue), value)
}

func (sm *SessionManager) Delete(ctx context.Context, key string) error {
	return sm.redisClient.Del(ctx, key).Err()
}

func IssueSession(c *gin.Context, token *oauth2.Token, userInfo map[string]interface{}) (string, error) {
	session := sessions.Default(c)

	sessionInfo := SessionInfo{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		Expiry:      token.Expiry,
		UserID:      userInfo["id"].(string),
		UserName:    userInfo["name"].(string),
		Email:       userInfo["email"].(string),
		// 设置其他必要的用户信息
	}

	session.Set("session_info", sessionInfo)

	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %v", err)
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	// 立即尝试检索会话信息
	checkSessionInfo := session.Get("session_info")
	if checkSessionInfo == nil {
		log.Println("Warning: Session info not found immediately after saving")
		return "", fmt.Errorf("session info not found after saving")
	}

	sessionID := session.ID()
	if sessionID == "" {
		log.Println("Warning: Session ID is empty")
		return "", fmt.Errorf("empty session ID")
	}

	log.Printf("Issuing new session. Session ID: %s", sessionID)
	log.Printf("Session info saved successfully: %+v", checkSessionInfo)

	return sessionID, nil
}
func init() {
	gob.Register(SessionInfo{})
}

func (sm SessionManager) GetSession(c *gin.Context) (*oauth2.Token, map[string]interface{}, error) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		return nil, nil, err
	}

	log.Printf("GET SESSION Retrieving session. Session ID: %s", sessionID)
	var sessionData map[string]interface{}
	err = sm.Get(c.Request.Context(), "session:"+sessionID, &sessionData)
	if err != nil {
		return nil, nil, err
	}

	if sessionData == nil {
		return nil, nil, errors.New("session data not found")
	}

	// 从 sessionData 中提取 token 信息
	tokenData, ok := sessionData["token"].(map[string]interface{})
	if !ok {
		return nil, nil, errors.New("invalid token data in session")
	}

	// 重构 oauth2.Token
	token := &oauth2.Token{
		AccessToken: tokenData["access_token"].(string),
		TokenType:   tokenData["token_type"].(string),
	}

	// 处理 expiry
	if expiryStr, ok := tokenData["expiry"].(string); ok {
		expiry, err := time.Parse(time.RFC3339, expiryStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid expiry format: %v", err)
		}
		token.Expiry = expiry
	} else {
		// 如果 expiry 不是字符串，可能是数字类型
		if expiryFloat, ok := tokenData["expiry"].(float64); ok {
			token.Expiry = time.Unix(int64(expiryFloat), 0)
		} else {
			return nil, nil, errors.New("expiry is neither string nor float64")
		}
	}

	// 从 sessionData 中提取 user_info
	userInfo, ok := sessionData["user_info"].(map[string]interface{})
	if !ok {
		return nil, nil, errors.New("invalid user info in session")
	}

	fmt.Println("name is ", userInfo["name"])

	return token, userInfo, nil
}

func SessionUsername(userInfo map[string]interface{}) string {
	if name, ok := userInfo["name"].(string); ok {
		return name
	}
	return ""
}

func DestroySession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	return session.Save()
}
