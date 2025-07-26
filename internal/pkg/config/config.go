package config

import (
	"github.com/spf13/viper"
)

// LoadConfig loads the configuration from the specified path and config name.
func LoadConfig(path string, configName string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}