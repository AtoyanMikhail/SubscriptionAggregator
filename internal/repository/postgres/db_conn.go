package postgres

import (
	"fmt"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/config"
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	log := logger.Global()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	log.Info("Connecting to PostgreSQL database",
		logger.String("host", cfg.Host),
		logger.String("port", cfg.Port),
		logger.String("database", cfg.DBName),
		logger.String("ssl_mode", cfg.SSLMode))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Error("Failed to connect to database",
			logger.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping database",
			logger.Error(err))
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL database")
	return db, nil
}
