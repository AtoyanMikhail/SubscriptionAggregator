package repository

import (
	"context"
	"time"
)

// SubscriptionsRepository defines the interface for database interaction
type SubscriptionsRepository interface {
	Create(ctx context.Context, subscription *Subscription) error
	GetSubscription(ctx context.Context, userID string, subscriptionID int) (*Subscription, error)
	UpdateSubscription(ctx context.Context, subscription *Subscription, userID string, subscriptionID int) error
	DeleteSubscription(ctx context.Context, userID string, subscriptionID int) error
	GetSubscriptionsByUserID(ctx context.Context, userID string) ([]*Subscription, error)
	GetSubscriptionsByPeriod(ctx context.Context, userID string, serviceNames []string, startDate, endDate time.Time) ([]*Subscription, error)
	Close() error
	RunMigrations(migrationsFilePath string) error
}