package app

import (
	"fmt"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/database"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/kafka"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/redisdb"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/handler"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/token"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log/slog"
)

type App struct {
	server *fiber.App
	db     *gorm.DB
	log    *slog.Logger
	port   string
}

func New(h *handler.FiberHandler, db *gorm.DB, log *slog.Logger, cfg ServerConfig) *App {
	server := fiber.New()
	server.Use(cors.New())
	h.RegisterRoutes(server)

	// Migrasi otomatis domain User dan VerificationCode
	db.AutoMigrate(&domain.User{}, &domain.VerificationCode{})

	return &App{
		server: server,
		db:     db,
		log:    log,
		port:   cfg.Port,
	}
}

func (a *App) Start() error {
	a.log.Info(fmt.Sprintf("Server auth dimulai di port %s", a.port))
	return a.server.Listen(":" + a.port)
}

// --- Konfigurasi Provider untuk Wire ---

type ServerConfig struct {
	Port string
}

func ProvideServerConfig() ServerConfig {
	return ServerConfig{Port: viper.GetString("app.port")}
}

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
	return kafka.ProducerConfig{
		BootstrapServers: viper.GetString("kafka.brokers"),
		Topic:            viper.GetString("kafka.topic"),
	}
}

func ProvideJWTConfig() token.JWTConfig {
	return token.JWTConfig{
		PrivateKeyPath: viper.GetString("jwt.private_key_path"),
		PublicKeyPath:  viper.GetString("jwt.public_key_path"),
	}
}

func ProvideRedisConfig() redisdb.Config {
	return redisdb.Config{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}
}

func ProvideRedisClient(cfg redisdb.Config) *redis.Client {
	return redisdb.NewClient(cfg)
}
