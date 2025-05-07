package config

import (
	"fmt"

	"github.com/CustomCloudStorage/databases"
	"github.com/CustomCloudStorage/services"
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Port     string `validate:"required"`
	Cors     CORSConfig
	Postgres databases.PostgresConfig `validate:"required"`
	Service  services.ServiceConfig   `validate:"required"`
}

type CORSConfig struct {
	AllowedOrigin string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("../config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}
