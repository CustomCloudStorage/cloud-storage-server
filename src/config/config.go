package config

import (
	"github.com/CustomCloudStorage/databases"
	"github.com/CustomCloudStorage/utils"
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Port       string `validate:"required"`
	Cors       CORSConfig
	Postgres   databases.PostgresConfig `validate:"required"`
	StorageDir string                   `validate:"required"`
	TmpUpload  string                   `validate:"required"`
}

type CORSConfig struct {
	AllowedOrigin string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("../config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, utils.ErrRead.Wrap(err, "failed to read config.yaml")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, utils.ErrUnmarshal.Wrap(err, "failed to unmarshal config")
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, utils.ErrValidation.Wrap(err, "config validation failed")
	}

	return &config, nil
}
