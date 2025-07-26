package app

import (
	"context"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/database"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/kafka"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/event"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/repository"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log/slog"
)

type App struct {
	consumer *event.KafkaConsumer
	log      *slog.Logger
	topic    string
}

func New(db *gorm.DB, consumer *event.KafkaConsumer, log *slog.Logger, cfg ServerConfig) *App {
	return &App{consumer: consumer, log: log, topic: cfg.Topic}
}

func (a *App) Start() error {
	a.log.Info("email service start")
	return a.consumer.Start(context.Background(), a.topic)
}

// config providers

type ServerConfig struct{ Topic string }

func ProvideServerConfig() ServerConfig { return ServerConfig{Topic: viper.GetString("kafka.topic")} }

func ProvideDBConfig() database.Config {
	return database.Config{
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.name"),
	}
}

func ProvideKafkaConfig() kafka.ProducerConfig {
	return kafka.ProducerConfig{BootstrapServers: viper.GetString("kafka.brokers"), Topic: viper.GetString("kafka.topic")}
}

func ProvideDB(dbCfg database.Config) (*gorm.DB, error) {
	return database.NewPostgresDB(dbCfg)
}

func ProvideConsumer(cfg kafka.ProducerConfig, repo *repository.PostgresLogRepository, log *slog.Logger) (*event.KafkaConsumer, func(), error) {
	c, err := kafka.NewConsumer(cfg.BootstrapServers, "email-service")
	if err != nil {
		return nil, nil, err
	}
	consumer := event.NewKafkaConsumer(c, repo, log)
	cleanup := func() { c.Close() }
	return consumer, cleanup, nil
}
