//go:build wireinject
// +build wireinject

package di

import (
	"log/slog"

	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/database"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/kafka"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/app"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/event"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/repository"
	"github.com/google/wire"
)

func InitializeApp(log *slog.Logger) (*app.App, func(), error) {
	wire.Build(
		app.New,
		repository.NewPostgresLogRepository,
		event.NewKafkaConsumer,
		database.NewPostgresDB,
		kafka.NewConsumer,
		app.ProvideDBConfig,
		app.ProvideKafkaConfig,
		app.ProvideServerConfig,
		app.ProvideDB,
		app.ProvideConsumer,
	)
	return &app.App{}, nil, nil
}
