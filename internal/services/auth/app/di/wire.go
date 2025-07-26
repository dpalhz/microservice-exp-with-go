//go:build wireinject
// +build wireinject

package di

import (
	"log/slog"

	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/database"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/kafka"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/app"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/event"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/handler"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/repository"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/session"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/token"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/usecase"
	"github.com/google/wire"
)

func InitializeApp(log *slog.Logger) (*app.App, func(), error) {
	wire.Build(
		// ⛓ BINDING INTERFACE TO IMPLEMENTATION — WAJIB!
		wire.Bind(new(usecase.UserRepository), new(*repository.PostgresUserRepository)),
		wire.Bind(new(usecase.EventProducer), new(*event.KafkaEventProducer)),
		wire.Bind(new(usecase.TokenGenerator), new(*token.JWTGenerator)),
		wire.Bind(new(usecase.SessionStore), new(*session.RedisSessionStore)),
		wire.Bind(new(usecase.VerificationRepository), new(*repository.PostgresVerificationRepository)),

		// ✅ PROVIDER CONSTRUCTORS
		app.New,
		handler.NewFiberHandler,
		usecase.NewUserUsecase,
		repository.NewPostgresUserRepository,
		repository.NewPostgresVerificationRepository,
		event.NewKafkaEventProducer,
		token.NewJWTGenerator,
		session.NewRedisSessionStore,
		app.ProvideRedisClient,
		database.NewPostgresDB,
		kafka.NewProducer,
		app.ProvideDBConfig,
		app.ProvideKafkaConfig,
		app.ProvideJWTConfig,
		app.ProvideRedisConfig,
		app.ProvideServerConfig,
	)
	return &app.App{}, nil, nil
}
