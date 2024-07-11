package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/gin-starter/internal/auth"
	"github.com/majiayu000/gin-starter/internal/types"
)

type AuthHandler struct {
	oauthManager   *auth.OAuthManager
	sessionManager types.SessionManager
}

func NewAuthHandler(om *auth.OAuthManager, sm types.SessionManager) *AuthHandler {
	return &AuthHandler{
		oauthManager:   om,
		sessionManager: sm,
	}
}

// HandleProfile handles the user profile request
func (h *AuthHandler) HandleProfile(c *gin.Context) {
	_, userInfo, err := h.sessionManager.GetSession(c)
	fmt.Print("user info is ", userInfo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	c.JSON(http.StatusOK, userInfo)
}

// HandleGoogleLogin initiates the Google OAuth login process
func (h *AuthHandler) HandleGoogleLogin(c *gin.Context) {
	h.Login(c)
}

// HandleGoogleCallback handles the Google OAuth callback
func (h *AuthHandler) HandleGoogleCallback(c *gin.Context) {
	h.Callback(c)
}

// HandleLogout logs out the user
func (h *AuthHandler) HandleLogout(c *gin.Context) {
	h.Logout(c)
}

func (h *AuthHandler) Login(c *gin.Context) {
	provider := c.Param("provider")
	providerInstance, err := h.oauthManager.GetProvider(provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	// 生成状态并存储在 Redis 中
	state := generateRandomState()
	err = h.sessionManager.Set(c.Request.Context(), "oauth_state:"+state, true, 10*time.Minute)
	if err != nil {
		log.Printf("Error setting state in Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize login"})
		return
	}
	log.Printf("Generated and stored state: %s", state)

	authURL := providerInstance.GetAuthURL(state)
	c.Redirect(http.StatusFound, authURL)
}

func (h *AuthHandler) Callback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	log.Printf("Received state in callback: %s", state)

	// 验证状态
	var stateExists bool
	err := h.sessionManager.Get(c.Request.Context(), "oauth_state:"+state, &stateExists)
	if err != nil || !stateExists {
		log.Printf("Invalid state: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	// 删除已使用的状态
	h.sessionManager.Delete(c.Request.Context(), "oauth_state:"+state)

	log.Printf("Received authorization code: %s", code)

	token, err := h.oauthManager.Exchange(provider, code)
	if err != nil {
		log.Printf("Exchange error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to exchange token: %v", err)})
		return
	}

	userInfo, err := h.oauthManager.GetUserInfo(provider, token)
	log.Printf("UserInfo: %+v", userInfo)
	if err != nil {
		log.Printf("User info error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// 创建会话并存储在 Redis 中
	sessionID := generateSessionID()
	sessionData := map[string]interface{}{
		"token":     token,
		"user_info": userInfo,
	}
	err = h.sessionManager.Set(c.Request.Context(), "session:"+sessionID, sessionData, 24*time.Hour)
	if err != nil {
		log.Printf("Session creation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// 设置 session cookie
	c.SetCookie(
		"session_id",
		sessionID,
		int(24*time.Hour.Seconds()), // 1 day in seconds
		"/",
		"",
		false,
		true,
	)

	log.Printf("Session created successfully. Session ID: %s", sessionID)

	// 设置 cookie

	// 验证会话是否正确保存
	var checkSessionInfo map[string]interface{}
	err = h.sessionManager.Get(c.Request.Context(), "session:"+sessionID, &checkSessionInfo)
	if err != nil {
		log.Println("Warning: Session info not found before redirect")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Session validation failed"})
		return
	}

	log.Printf("Session verified successfully before redirect. Session ID: %s", sessionID)

	c.Redirect(http.StatusFound, "/")
}

func generateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if err := auth.DestroySession(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to destroy session"})
		return
	}
	c.Redirect(http.StatusFound, "/")
}
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	_, userInfo, err := h.sessionManager.GetSession(c)
	fmt.Println("get session", userInfo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	username := auth.SessionUsername(userInfo)
	c.JSON(http.StatusOK, gin.H{"username": username})
}
