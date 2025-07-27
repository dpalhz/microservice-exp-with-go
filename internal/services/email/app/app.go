package app

import (
	"log/slog"

	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/database"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/kafka"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/domain"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/event"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ServerConfig struct {
	Port  string
}

type App struct {
	server 	 *fiber.App
	db     	 *gorm.DB
	consumer *event.KafkaConsumer
	log      *slog.Logger
	port      string
	
}

func New(consumer *event.KafkaConsumer, db *gorm.DB, log *slog.Logger, cfg ServerConfig) *App {
	server := fiber.New()
	server.Use(cors.New())
	server.Use(fiber.New().Use(cors.New()))

	db.AutoMigrate(&domain.EmailLog{})

	return &App{
		server:   server,
		db:       db,	
		consumer: consumer,	
		log:      log,
	}
}


func (a *App) Start(cfg ServerConfig) error {
	a.log.Info("Email service started at " + a.port)
	return a.server.Listen(":" + a.port) 
}



// config providers
func ProvideServerConfig() ServerConfig {
	return ServerConfig{
		Port:  viper.GetString("app.port"),
	}
}


func ProvideDB(dbCfg database.Config) (*gorm.DB, error) {
	return database.NewPostgresDB(dbCfg)
}

func ProvideConsumer(cfg kafka.ConsumerConfig, repo *repository.PostgresLogRepository, log *slog.Logger) (*event.KafkaConsumer, func(), error) {
	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, nil, err
	}
	consumer := event.NewKafkaConsumer(c, repo, log)
	cleanup := func() { c.Close() }
	return consumer, cleanup, nil
}
