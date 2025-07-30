package postgres

import "errors"

var (
	// Create subscription errors
	ErrCreateSubscriptionFailed = errors.New("failed to create subscription")

	// Get subscription errors
	ErrGetSubscriptionFailed = errors.New("failed to get subscription")
	ErrSubscriptionNotFound  = errors.New("subscription not found")

	// Update subscription errors
	ErrUpdateSubscriptionFailed      = errors.New("failed to update subscription")
	ErrGetRowsAffectedFailed         = errors.New("failed to get rows affected")
	ErrSubscriptionNotFoundForUpdate = errors.New("subscription not found")

	// Delete subscription errors
	ErrDeleteSubscriptionFailed        = errors.New("failed to delete subscription")
	ErrSubscriptionNotFoundForDeletion = errors.New("subscription not found")

	// Get subscriptions by user ID errors
	ErrGetSubscriptionsByUserIDFailed = errors.New("failed to get subscriptions")

	// Get subscriptions by period errors
	ErrGetSubscriptionsByPeriodFailed = errors.New("failed to get subscriptions by period")

	// Migration errors
	ErrCreateMigrationDriverFailed   = errors.New("failed to create migration driver")
	ErrCreateMigrationInstanceFailed = errors.New("failed to create migration instance")
	ErrRunMigrationsFailed           = errors.New("failed to run migrations")
)
