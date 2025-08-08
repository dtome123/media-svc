package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server   `mapstructure:"server"`
	DB       DB       `mapstructure:"db"`
	S3       S3       `mapstructure:"s3"`
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
}

type RabbitMQ struct {
	DSN   string `mapstructure:"dsn"`
	Queue string `mapstructure:"queue"`
}

type S3 struct {
	Endpoint     string `mapstructure:"endpoint"`
	Bucket       string `mapstructure:"bucket"`
	StreamBucket string `mapstructure:"stream_bucket"`
	AccessKey    string `mapstructure:"access_key"`
	SecretKey    string `mapstructure:"secret_key"`
	Secure       bool   `mapstructure:"secure"`
	Region       string `mapstructure:"region"`
}

type Server struct {
	GrpcPort string `mapstructure:"grpc_port"`
	HttpPort string `mapstructure:"http_port"`
}

type DB struct {
	Mongo struct {
		DSN      string `mapstructure:"dsn"`
		Database string `mapstructure:"database"`
	} `mapstructure:"mongo"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	v := viper.NewWithOptions()
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	v.SetConfigFile("config/config.yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
