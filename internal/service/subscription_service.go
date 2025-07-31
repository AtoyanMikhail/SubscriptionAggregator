package service

import (
	"context"
	"fmt"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/logger"
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/repository"
	"github.com/go-playground/validator/v10"
)

type subscriptionService struct {
	repo      repository.SubscriptionsRepository
	log       logger.Logger
	validator *validator.Validate
}

// NewSubscriptionService creates a new instance of subscription service
func NewSubscriptionService(repo repository.SubscriptionsRepository) SubscriptionService {
	return &subscriptionService{
		repo:      repo,
		log:       logger.Global(),
		validator: validator.New(),
	}
}

// CreateSubscription creates a new subscription
func (s *subscriptionService) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*repository.Subscription, error) {

	s.log.Info("creating new subscription",
		logger.String("user_id", req.UserID),
		logger.String("service_name", req.ServiceName))

	// Validate request
	if err := s.validator.Struct(req); err != nil {
		s.log.Error("subscription creation validation failed",
			logger.Error(err),
			logger.String("user_id", req.UserID))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert to model
	subscription, err := req.ToSubscriptionModel()
	if err != nil {
		s.log.Error("failed to convert request to subscription model",
			logger.Error(err),
			logger.String("user_id", req.UserID))
		return nil, err
	}

	// Create subscription
	if err := s.repo.Create(ctx, subscription); err != nil {
		s.log.Error("failed to create subscription in repository",
			logger.Error(err),
			logger.String("user_id", req.UserID))
		return nil, ErrInternalServer
	}

	s.log.Info("subscription created successfully",
		logger.Int("subscription_id", subscription.ID),
		logger.String("user_id", req.UserID))

	return subscription, nil
}

// GetSubscription retrieves a specific subscription
func (s *subscriptionService) GetSubscription(ctx context.Context, userID string, subscriptionID int) (*repository.Subscription, error) {
	s.log.Debug("getting subscription",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	// Validate user ID
	if err := s.validator.Var(userID, "required,uuid4"); err != nil {
		s.log.Error("invalid user ID format",
			logger.Error(err),
			logger.String("user_id", userID))
		return nil, ErrInvalidUserID
	}

	// Validate subscription ID
	if subscriptionID <= 0 {
		s.log.Error("invalid subscription ID",
			logger.Int("subscription_id", subscriptionID))
		return nil, ErrInvalidSubscriptionID
	}

	subscription, err := s.repo.GetSubscription(ctx, userID, subscriptionID)
	if err != nil {
		s.log.Error("failed to get subscription from repository",
			logger.Error(err),
			logger.String("user_id", userID),
			logger.Int("subscription_id", subscriptionID))
		return nil, ErrSubscriptionNotFound
	}

	s.log.Debug("subscription retrieved successfully",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	return subscription, nil
}

// UpdateSubscription updates an existing subscription
func (s *subscriptionService) UpdateSubscription(ctx context.Context, userID string, subscriptionID int, req *UpdateSubscriptionRequest) (*repository.Subscription, error) {
	s.log.Info("updating subscription",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	// Validate user ID
	if err := s.validator.Var(userID, "required,uuid4"); err != nil {
		s.log.Error("invalid user ID format",
			logger.Error(err),
			logger.String("user_id", userID))
		return nil, ErrInvalidUserID
	}

	// Validate subscription ID
	if subscriptionID <= 0 {
		s.log.Error("invalid subscription ID",
			logger.Int("subscription_id", subscriptionID))
		return nil, ErrInvalidSubscriptionID
	}

	// Validate request
	if err := s.validator.Struct(req); err != nil {
		s.log.Error("subscription update validation failed",
			logger.Error(err),
			logger.String("user_id", userID),
			logger.Int("subscription_id", subscriptionID))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert to model
	subscription, err := req.ToSubscriptionModel()
	if err != nil {
		s.log.Error("failed to convert request to subscription model",
			logger.Error(err),
			logger.String("user_id", userID))
		return nil, err
	}

	// Set ID for the model
	subscription.ID = subscriptionID
	subscription.UserID = userID

	// Update subscription
	if err := s.repo.UpdateSubscription(ctx, subscription, userID, subscriptionID); err != nil {
		s.log.Error("failed to update subscription in repository",
			logger.Error(err),
			logger.String("user_id", userID),
			logger.Int("subscription_id", subscriptionID))
		return nil, ErrSubscriptionNotFound
	}

	s.log.Info("subscription updated successfully",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	return subscription, nil
}

// DeleteSubscription deletes a subscription
func (s *subscriptionService) DeleteSubscription(ctx context.Context, userID string, subscriptionID int) error {
	s.log.Info("deleting subscription",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	// Validate user ID
	if err := s.validator.Var(userID, "required,uuid4"); err != nil {
		s.log.Error("invalid user ID format",
			logger.Error(err),
			logger.String("user_id", userID))
		return ErrInvalidUserID
	}

	// Validate subscription ID
	if subscriptionID <= 0 {
		s.log.Error("invalid subscription ID",
			logger.Int("subscription_id", subscriptionID))
		return ErrInvalidSubscriptionID
	}

	if err := s.repo.DeleteSubscription(ctx, userID, subscriptionID); err != nil {
		s.log.Error("failed to delete subscription from repository",
			logger.Error(err),
			logger.String("user_id", userID),
			logger.Int("subscription_id", subscriptionID))
		return ErrSubscriptionNotFound
	}

	s.log.Info("subscription deleted successfully",
		logger.String("user_id", userID),
		logger.Int("subscription_id", subscriptionID))

	return nil
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (s *subscriptionService) GetUserSubscriptions(ctx context.Context, userID string) ([]*repository.Subscription, error) {
	s.log.Debug("getting user subscriptions",
		logger.String("user_id", userID))

	// Validate user ID
	if err := s.validator.Var(userID, "required,uuid4"); err != nil {
		s.log.Error("invalid user ID format",
			logger.Error(err),
			logger.String("user_id", userID))
		return nil, ErrInvalidUserID
	}

	subscriptions, err := s.repo.GetSubscriptionsByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user subscriptions from repository",
			logger.Error(err),
			logger.String("user_id", userID))
		return nil, ErrInternalServer
	}

	s.log.Debug("user subscriptions retrieved successfully",
		logger.String("user_id", userID),
		logger.Int("count", len(subscriptions)))

	return subscriptions, nil
}

// CalculateTotalCost calculates total cost of subscriptions for a period
func (s *subscriptionService) CalculateTotalCost(ctx context.Context, req *GetCostRequest) (*CostResponse, error) {
	s.log.Info("calculating total subscription cost",
		logger.String("user_id", req.UserID),
		logger.String("start_date", req.StartDate),
		logger.String("end_date", req.EndDate))

	// Validate request
	if err := s.validator.Struct(req); err != nil {
		s.log.Error("cost calculation validation failed",
			logger.Error(err),
			logger.String("user_id", req.UserID))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Parse dates
	startDate, err := ParseMonthYear(req.StartDate)
	if err != nil {
		s.log.Error("failed to parse start date",
			logger.Error(err),
			logger.String("start_date", req.StartDate))
		return nil, err
	}

	endDate, err := ParseMonthYear(req.EndDate)
	if err != nil {
		s.log.Error("failed to parse end date",
			logger.Error(err),
			logger.String("end_date", req.EndDate))
		return nil, err
	}

	// Set end date to last day of the month
	endDate = GetLastDayOfMonth(endDate)

	// Validate date range
	if endDate.Before(startDate) {
		s.log.Error("end date is before start date",
			logger.String("start_date", req.StartDate),
			logger.String("end_date", req.EndDate))
		return nil, ErrInvalidDateRange
	}

	// Get subscriptions for the period
	subscriptions, err := s.repo.GetSubscriptionsByPeriod(ctx, req.UserID, req.ServiceNames, startDate, endDate)
	if err != nil {
		s.log.Error("failed to get subscriptions for period",
			logger.Error(err),
			logger.String("user_id", req.UserID))
		return nil, ErrInternalServer
	}

	// Calculate costs
	var totalCost int
	breakdown := make([]SubscriptionCostBreakdown, 0, len(subscriptions))

	for _, sub := range subscriptions {
		monthsCount := CalculateSubscriptionMonthsInPeriod(&sub.StartDate, sub.EndDate, startDate, endDate)
		if monthsCount <= 0 {
			continue
		}

		subTotalCost := sub.Price * monthsCount
		totalCost += subTotalCost

		breakdown = append(breakdown, SubscriptionCostBreakdown{
			SubscriptionID: sub.ID,
			ServiceName:    sub.ServiceName,
			MonthlyPrice:   sub.Price,
			MonthsCount:    monthsCount,
			TotalCost:      subTotalCost,
		})

		s.log.Debug("calculated cost for subscription",
			logger.Int("subscription_id", sub.ID),
			logger.String("service_name", sub.ServiceName),
			logger.Int("months_count", monthsCount),
			logger.Int("total_cost", subTotalCost))
	}

	response := &CostResponse{
		UserID:    req.UserID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		TotalCost: totalCost,
		Breakdown: breakdown,
	}

	s.log.Info("total subscription cost calculated successfully",
		logger.String("user_id", req.UserID),
		logger.Int("total_cost", totalCost),
		logger.Int("subscriptions_count", len(breakdown)))

	return response, nil
}
