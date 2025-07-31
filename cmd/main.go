package main

import (
	"fmt"
	"os"

	_ "github.com/AtoyanMikhail/SubscribtionAggregation/docs" // swagger docs
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/config"
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/handlers"
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/logger"
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/repository/postgres"
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/service"
)

// @title Subscription Management API
// @version 1.0
// @description REST API for managing user subscriptions

// @host localhost:8080
// @BasePath /

// @schemes http https

func main() {
	// Initialize logger
	logger.Initialize(os.Stdout)
	log := logger.Global()

	// Load configuration
	cfg, err := config.GetConfig("configs/app/config_local.yaml")
	if err != nil {
		log.Fatal("failed to load configuration", logger.Error(err))
	}

	fmt.Println(cfg)

	// Connect to database
	db, err := postgres.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", logger.Error(err))
	}
	defer db.Close()

	// Initialize repository
	subscriptionRepo := postgres.NewSubscriptionsRepository(db)

	// Run migrations
	if err := subscriptionRepo.RunMigrations("migrations"); err != nil {
		log.Fatal("failed to run migrations", logger.Error(err))
	}

	// Initialize service
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	// Setup router
	router := handlers.SetupRouter(subscriptionService)

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Info("starting HTTP server",
		logger.String("address", serverAddr))

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("failed to start server", logger.Error(err))
	}
}
