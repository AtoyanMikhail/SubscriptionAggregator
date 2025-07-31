package config

import (
	"sync"
)

var (
	globalConfig Config
	initOnce     sync.Once
)

type Config struct {
	Server   ServerConfig   `yaml:"server" envPrefix:"SERVER_" validate:"required"`
	Database DatabaseConfig `yaml:"database" envPrefix:"DB_" validate:"required"`
}

type ServerConfig struct {
	Port         string   `yaml:"port" env:"PORT" validate:"required,numeric"`
	Host         string   `yaml:"host" env:"HOST" validate:"required,hostname|ip"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"HOST" validate:"required,hostname|ip"`
	Port     string `yaml:"port" env:"PORT" validate:"required,numeric"`
	User     string `yaml:"user" env:"USER" validate:"required"`
	Password string `yaml:"password" env:"PASSWORD" validate:"required"`
	DBName   string `yaml:"db_name" env:"NAME" validate:"required"`
	SSLMode  string `yaml:"ssl_mode" env:"SSL_MODE" validate:"required,oneof=disable require verify-ca verify-full"`
}
