package models

import (
	"context"
)

type SubscriptionsRepository interface {
	Create(ctx context.Context, subscription *Subscription) error
	GetSubscription(ctx context.Context, userID, subscriptionID string) (*Subscription, error)
	UpdateSubscription(ctx context.Context, userID, subscriptionID string) error
	DeleteSubscription(ctx context.Context, userID, subscriptionID string) error
	GetSubscriptionsByUserID(ctx context.Context, userID string) ([]*Subscription, error)
	Close() error
	RunMigrations(migrationsFilePath string) error
}
