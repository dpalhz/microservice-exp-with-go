package main

import (
	"log/slog"

	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/config"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/logger"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/app"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/app/di"
)

func main() {
	log := logger.New()
	if err := config.LoadConfig("./configs", "email"); err != nil {
		log.Error("failed load config", slog.String("error", err.Error()))
		return
	}
	application, cleanup, err := di.InitializeApp(log)
	if err != nil {
		log.Error("init app", slog.String("error", err.Error()))
		return
	}
	defer cleanup()
	if err := application.Start(); err != nil {
		log.Error("start", slog.String("error", err.Error()))
	}
}
