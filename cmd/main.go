// cmd/main.go
package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	config "github.com/majiayu000/gin-starter/configs"

	"github.com/majiayu000/gin-starter/internal/auth"
	"github.com/majiayu000/gin-starter/internal/auth/oauth"
	"github.com/majiayu000/gin-starter/internal/database"
	"github.com/majiayu000/gin-starter/internal/handlers"
	"github.com/majiayu000/gin-starter/internal/router"
)

func main() {

	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("cfg is", cfg.OAuth.Apple)
	keyPath := filepath.Join(".", cfg.OAuth.Apple.KeyPath)
	fmt.Println("key path is", keyPath)
	// fmt.Printf("Loaded config: %+v\n", cfg)
	fmt.Printf("MySQL Host: %s\n", cfg.DB.MySQL.Host)
	// fmt.Printf("MySQL Port: %d\n", config.MySQL.Port)
	oauthConfig := map[string]map[string]string{
		"google": {
			"client_id":     cfg.OAuth.Google.ClientID,
			"client_secret": cfg.OAuth.Google.ClientSecret,
			"redirect_url":  cfg.OAuth.Google.RedirectURL,
		},
		"apple": {
			"client_id":    cfg.OAuth.Apple.ClientID,
			"team_id":      cfg.OAuth.Apple.TeamID,
			"key_id":       cfg.OAuth.Apple.KeyID,
			"private_key":  cfg.OAuth.Apple.KeyPath,
			"redirect_url": cfg.OAuth.Apple.RedirectURL,
		},
	}
	fmt.Println("config is ", cfg.OAuth.Google.RedirectURL)
	oauthManager := auth.NewOAuthManager()

	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "123456", []byte("secret"))
	store.Options(sessions.Options{
		MaxAge:   3600 * 24, // 24 小时
		HttpOnly: true,
		Secure:   true, // 如果使用 HTTPS
	})
	sessionManager := auth.NewSessionManager("localhost:6379", "123456", 0)
	googleProvider, _ := oauth.NewGoogleProvider(
		oauthConfig["google"],
		sessionManager,
	)

	oauthManager.AddProvider("google", googleProvider)
	authHandler := handlers.NewAuthHandler(oauthManager, sessionManager)

	// 初始化数据库连接
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置路由
	r := router.SetupRouter(db)

	r.Use(sessions.Sessions("mysession", store))

	r.GET("/", authHandler.HandleProfile)
	r.GET("/auth/:provider/login", authHandler.HandleGoogleLogin)
	r.GET("/auth/:provider/callback", authHandler.HandleGoogleCallback)
	r.POST("/logout", authHandler.HandleLogout)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// log.Printf("Starting Server listening on localhost:8080")
	// if err := r.Run(":8080"); err != nil {
	// 	log.Fatal("Run: ", err)
	// }
}
