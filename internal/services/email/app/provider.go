package app

import (
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/database"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/kafka"
	"github.com/spf13/viper"
)

func ProvideDBConfig() database.Config {
	return database.Config{
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.name"),
	}
}

func ProvideKafkaConsumerConfig() kafka.ConsumerConfig {
	return kafka.ConsumerConfig{
		BootstrapServers: viper.GetString("kafka.brokers"),
		GroupID:          viper.GetString("kafka.group_id"),
	}
}
