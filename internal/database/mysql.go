package database

import (
	"fmt"

	config "github.com/majiayu000/gin-starter/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.MySQL.User,
		cfg.DB.MySQL.Password,
		cfg.DB.MySQL.Host,
		cfg.DB.MySQL.Port,
		cfg.DB.MySQL.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
