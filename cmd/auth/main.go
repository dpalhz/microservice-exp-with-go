package main

import (
	"log/slog"

	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/config"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/logger"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/app/di"
)

func main() {
	log := logger.New()

	if err := config.LoadConfig("./configs", "auth"); err!= nil {
		log.Error("Gagal memuat konfigurasi", slog.String("error", err.Error()))
		return
	}

	app, cleanup, err := di.InitializeApp(log)
	if err!= nil {
		log.Error("Gagal menginisialisasi aplikasi", slog.String("error", err.Error()))
		return
	}
	defer cleanup()

	if err := app.Start(); err!= nil {
		log.Error("Gagal memulai server", slog.String("error", err.Error()))
	}
}