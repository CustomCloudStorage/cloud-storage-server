package config

import (
	"github.com/CustomCloudStorage/databases"
	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	Cors     CORSConfig
	Postgres databases.PostgresConfig
}

type CORSConfig struct {
	AllowedOrigin string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("../config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
