package config

import (
	"os"

	"log"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

// GetConfig sets default values to the Config struct, then tries to override them with a .json config file(configPath),
// and finally overrides values from environment variables on the first usage. Then, it returns a pointer to the global config instance.
func GetConfig(configPath string) (*Config, error) {
	initOnce.Do(func() {
		setDefaults(&globalConfig)

		// Overriding default values with json if there is a valid config
		if err := loadFromYAML(configPath, &globalConfig); err != nil {
			log.Printf("failed to load config from json: %s\n", err.Error())
		}

		// Overriding json values with env
		loadFromEnv(&globalConfig)

		if err := validate(&globalConfig); err != nil {
			log.Fatalf("config validation failed: %s", err.Error())
		}
	})

	return &globalConfig, nil
}

func setDefaults(cfg *Config) {
	cfg.Server = ServerConfig{
		Port: "8080",
		Host: "0.0.0.0",
	}

	cfg.Database = DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		DBName:   "subscription_aggregator",
		SSLMode:  "disable",
	}
}

func loadFromYAML(path string, cfg *Config) error {
	configPath := path
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

func loadFromEnv(cfg *Config) {
	_ = env.Parse(cfg)
}

func validate(cfg *Config) error {
	validate := validator.New()

	return validate.Struct(cfg)
}
