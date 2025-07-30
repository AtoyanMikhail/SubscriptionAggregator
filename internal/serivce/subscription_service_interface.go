package service

import (
	"context"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/repository"
)

// SubscriptionService defines the interface for subscription business logic
type SubscriptionService interface {
	// CRUD operations
	CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*repository.Subscription, error)
	GetSubscription(ctx context.Context, userID string, subscriptionID int) (*repository.Subscription, error)
	UpdateSubscription(ctx context.Context, userID string, subscriptionID int, req *UpdateSubscriptionRequest) (*repository.Subscription, error)
	DeleteSubscription(ctx context.Context, userID string, subscriptionID int) error
	GetUserSubscriptions(ctx context.Context, userID string) ([]*repository.Subscription, error)

	// Cost calculation
	CalculateTotalCost(ctx context.Context, req *GetCostRequest) (*CostResponse, error)
}
