package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

var cfg *Config

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	cfg = &config
	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}
