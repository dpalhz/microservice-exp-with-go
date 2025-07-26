package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func NewPostgresDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err!= nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}