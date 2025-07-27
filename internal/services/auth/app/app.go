package app

import (
	"fmt"
	"log/slog"

	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

type ServerConfig struct {
	Port string
}


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

	db.AutoMigrate(&domain.User{}, &domain.VerificationCode{})

	return &App{
		server: server,
		db:     db,
		log:    log,
		port:   cfg.Port,
	}
}

func (a *App) Start() error {
	a.log.Info(fmt.Sprintf("server started at %s", a.port))
	return a.server.Listen(":" + a.port)
}


