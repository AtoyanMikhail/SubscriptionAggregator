package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/logger"
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type subscriptionsRepository struct {
	db *sqlx.DB
}

// NewSubscriptionsRepository creates a new instance of PostgreSQL subscriptions repository
func NewSubscriptionsRepository(db *sqlx.DB) repository.SubscriptionsRepository {
	return &subscriptionsRepository{
		db: db,
	}
}

// Create inserts a new subscription into the database
func (r *subscriptionsRepository) Create(ctx context.Context, subscription *repository.Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	log := logger.Global()
	log.Debug("Creating subscription",
		logger.String("user_id", subscription.UserID),
		logger.String("service_name", subscription.ServiceName),
		logger.Int("price", subscription.Price))

	err := r.db.QueryRowContext(ctx, query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate).Scan(&subscription.ID)

	if err != nil {
		log.Error("Failed to create subscription",
			logger.Error(err),
			logger.String("user_id", subscription.UserID))
		return ErrCreateSubscriptionFailed
	}

	log.Info("Subscription created successfully",
		logger.Int("subscription_id", subscription.ID),
		logger.String("user_id", subscription.UserID))

	return nil
}

// GetSubscription retrieves a specific subscription by user ID and subscription ID
func (r *subscriptionsRepository) GetSubscription(ctx context.Context, userID string, subscriptionID int) (*repository.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1 AND id = $2`

	log := logger.Global()
	log.Debug("Getting subscription",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	subscription := &repository.Subscription{}
	err := r.db.GetContext(ctx, subscription, query, userID, subscriptionID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn("Subscription not found",
				logger.String("user_id", userID),
				logger.Int("subscription_id", subscriptionID))
			return nil, ErrSubscriptionNotFound
		}
		log.Error("Failed to get subscription",
			logger.Error(err),
			logger.String("user_id", userID),
			logger.Int("subscription_id", subscriptionID))
		return nil, ErrGetSubscriptionFailed
	}

	log.Debug("Subscription retrieved successfully",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	return subscription, nil
}

// UpdateSubscription updates an existing subscription
func (r *subscriptionsRepository) UpdateSubscription(ctx context.Context, subscription *repository.Subscription, userID string, subscriptionID int) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4
		WHERE user_id = $5 AND id = $6`

	log := logger.Global()
	log.Debug("Updating subscription",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	result, err := r.db.ExecContext(ctx, query,
		subscription.ServiceName,
		subscription.Price,
		subscription.StartDate,
		subscription.EndDate,
		userID,
		subscriptionID)

	if err != nil {
		log.Error("Failed to update subscription",
			logger.Error(err),
			logger.String("user_id", userID),
			logger.Int("subscription_id", subscriptionID))
		return ErrUpdateSubscriptionFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Failed to get rows affected",
			logger.Error(err))
		return ErrGetRowsAffectedFailed
	}

	if rowsAffected == 0 {
		log.Warn("Subscription not found for update",
			logger.String("user_id", userID),
			logger.Int("subscription_id", subscriptionID))
		return ErrSubscriptionNotFoundForUpdate
	}

	log.Info("Subscription updated successfully",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	return nil
}

// DeleteSubscription removes a subscription from the database
func (r *subscriptionsRepository) DeleteSubscription(ctx context.Context, userID string, subscriptionID int) error {
	query := `DELETE FROM subscriptions WHERE user_id = $1 AND id = $2`

	log := logger.Global()
	log.Debug("Deleting subscription",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	result, err := r.db.ExecContext(ctx, query, userID, subscriptionID)
	if err != nil {
		log.Error("Failed to delete subscription",
			logger.Error(err),
			logger.String("user_id", userID),
			logger.Int("subscription_id", subscriptionID))
		return ErrDeleteSubscriptionFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Failed to get rows affected",
			logger.Error(err))
		return ErrGetRowsAffectedFailed
	}

	if rowsAffected == 0 {
		log.Warn("Subscription not found for deletion",
			logger.String("user_id", userID),
			logger.Int("subscription_id", subscriptionID))
		return ErrSubscriptionNotFoundForDeletion
	}

	log.Info("Subscription deleted successfully",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	return nil
}

// GetSubscriptionsByUserID retrieves all subscriptions for a specific user
func (r *subscriptionsRepository) GetSubscriptionsByUserID(ctx context.Context, userID string) ([]*repository.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY start_date DESC`

	log := logger.Global()
	log.Debug("Getting subscriptions by user ID",
		logger.String("user_id", userID))

	subscriptions := []*repository.Subscription{}
	err := r.db.SelectContext(ctx, &subscriptions, query, userID)

	if err != nil {
		log.Error("Failed to get subscriptions by user ID",
			logger.Error(err),
			logger.String("user_id", userID))
		return nil, ErrGetSubscriptionsByUserIDFailed
	}

	log.Debug("Subscriptions retrieved successfully",
		logger.String("user_id", userID),
		logger.Int("count", len(subscriptions)))

	return subscriptions, nil
}

// GetSubscriptionsByPeriod retrieves subscriptions for a user within a time period with optional service name filtering
func (r *subscriptionsRepository) GetSubscriptionsByPeriod(ctx context.Context, userID string, serviceNames []string, startDate, endDate time.Time) ([]*repository.Subscription, error) {
	log := logger.Global()
	log.Debug("Getting subscriptions by period",
		logger.String("user_id", userID),
		logger.Any("service_names", serviceNames),
		logger.Any("start_date", startDate),
		logger.Any("end_date", endDate))

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
		AND start_date <= $3
		AND (end_date IS NULL OR end_date >= $2)`)

	args := []interface{}{userID, startDate, endDate}
	argIndex := 4

	// Add service name filtering if provided
	if len(serviceNames) > 0 {
		placeholders := make([]string, len(serviceNames))
		for i, serviceName := range serviceNames {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, serviceName)
			argIndex++
		}
		queryBuilder.WriteString(fmt.Sprintf(" AND service_name IN (%s)", strings.Join(placeholders, ",")))
	}

	queryBuilder.WriteString(" ORDER BY start_date DESC")

	subscriptions := []*repository.Subscription{}
	err := r.db.SelectContext(ctx, &subscriptions, queryBuilder.String(), args...)

	if err != nil {
		log.Error("Failed to get subscriptions by period",
			logger.Error(err),
			logger.String("user_id", userID))
		return nil, ErrGetSubscriptionsByPeriodFailed
	}

	log.Debug("Subscriptions by period retrieved successfully",
		logger.String("user_id", userID),
		logger.Int("count", len(subscriptions)))

	return subscriptions, nil
}

// Close closes the database connection
func (r *subscriptionsRepository) Close() error {
	log := logger.Global()
	log.Info("Closing database connection")
	return r.db.Close()
}

// RunMigrations runs database migrations
func (r *subscriptionsRepository) RunMigrations(migrationsFilePath string) error {
	log := logger.Global()
	log.Info("Running database migrations",
		logger.String("migrations_path", migrationsFilePath))

	driver, err := postgres.WithInstance(r.db.DB, &postgres.Config{})
	if err != nil {
		log.Error("Failed to create migration driver",
			logger.Error(err))
		return ErrCreateMigrationDriverFailed
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsFilePath),
		"postgres",
		driver,
	)
	if err != nil {
		log.Error("Failed to create migration instance",
			logger.Error(err))
		return ErrCreateMigrationInstanceFailed
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error("Failed to run migrations",
			logger.Error(err))
		return ErrRunMigrationsFailed
	}

	log.Info("Database migrations completed successfully")
	return nil
}