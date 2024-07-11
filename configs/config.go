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
	fmt.Println(config)

	return &config, nil
}
