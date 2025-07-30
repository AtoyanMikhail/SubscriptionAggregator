package config

import (
	"sync"
)

var (
	globalConfig Config
	initOnce     sync.Once
)

type Config struct {
	Server   ServerConfig   `json:"server" envPrefix:"SERVER_" validate:"required"`
	Database DatabaseConfig `json:"database" envPrefix:"DB_" validate:"required"`
}

type ServerConfig struct {
	Port         string   `json:"port" env:"PORT" validate:"required,numeric"`
	Host         string   `json:"host" env:"HOST" validate:"required,hostname|ip"`
	ReadTimeout  Duration `json:"read_timeout" env:"READ_TIMEOUT" validate:"required,duration_gt0"`
	WriteTimeout Duration `json:"write_timeout" env:"WRITE_TIMEOUT" validate:"required,duration_gt0"`
}

type DatabaseConfig struct {
	Host     string `json:"host" env:"HOST" validate:"required,hostname|ip"`
	Port     string `json:"port" env:"PORT" validate:"required,numeric"`
	User     string `json:"user" env:"USER" validate:"required"`
	Password string `json:"password" env:"PASSWORD" validate:"required"`
	DBName   string `json:"db_name" env:"NAME" validate:"required"`
	SSLMode  string `json:"ssl_mode" env:"SSL_MODE" validate:"required,oneof=disable require verify-ca verify-full"`
}
