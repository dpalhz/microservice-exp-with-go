package app

import (
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/database"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/kafka"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/redisdb"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/token"
	"github.com/spf13/viper"
)

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


