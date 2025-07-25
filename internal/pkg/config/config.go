package config

import (
	"github.com/spf13/viper"
)

// LoadConfig memuat konfigurasi dari file dan variabel lingkungan.
func LoadConfig(path string, configName string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}