// configs/config.go
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	OAuth struct {
		Google struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
			RedirectURL  string `mapstructure:"redirect_url"`
		} `mapstructure:"google"`
		Apple struct {
			ClientID    string `mapstructure:"client_id"`
			TeamID      string `mapstructure:"team_id"`
			KeyID       string `mapstructure:"key_id"`
			KeyPath     string `mapstructure:"key_path"`
			RedirectURL string `mapstructure:"redirect_url"`
		} `mapstructure:"apple"`
	} `mapstructure:"oauth"`
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	DB struct {
		MySQL struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			Name     string `mapstructure:"name"`
			Prefix   string `mapstructure:"prefix"`
		} `mapstructure:"mysql"`
		Redis struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			Name     string `mapstructure:"name"`
			Prefix   string `mapstructure:"prefix"`
		} `mapstructure:"redis"`
	} `mapstructure:"db"`
	LLM struct {
		AzureOpenAI struct {
			APIKey              string `mapstructure:"api_key"`
			Endpoint            string `mapstructure:"endpoint"`
			DeploymentName      string `mapstructure:"deployment_name"`
			SystemPromptZh      string `mapstructure:"system_prompt_zh"`
			SystemPromptEn      string `mapstructure:"system_prompt_en"`
			SystemPromptDefault string `mapstructure:"system_prompt_default"`
		} `mapstructure:"azure_openai"`
	} `mapstructure:"llm"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 从环境变量覆盖特定值
	if envClientID := viper.GetString("GOOGLE_CLIENT_ID"); envClientID != "" {
		config.OAuth.Google.ClientID = envClientID
	}
	if envClientSecret := viper.GetString("GOOGLE_CLIENT_SECRET"); envClientSecret != "" {
		config.OAuth.Google.ClientSecret = envClientSecret
	}

	// 从环境变量覆盖MySQL配置

	fmt.Printf("Loaded config: %+v\n", config)
	fmt.Printf("APIKey: %s\n", config.LLM.AzureOpenAI.APIKey)

	return &config, nil
}
