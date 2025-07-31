package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockSubscriptionsRepository is a mock implementation of SubscriptionsRepository
type MockSubscriptionsRepository struct {
	mock.Mock
}

func (m *MockSubscriptionsRepository) Create(ctx context.Context, subscription *repository.Subscription) error {
	args := m.Called(ctx, subscription)
	// Simulate setting ID after creation
	if args.Error(0) == nil {
		subscription.ID = 1
	}
	return args.Error(0)
}

func (m *MockSubscriptionsRepository) GetSubscription(ctx context.Context, userID string, subscriptionID int) (*repository.Subscription, error) {
	args := m.Called(ctx, userID, subscriptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.Subscription), args.Error(1)
}

func (m *MockSubscriptionsRepository) UpdateSubscription(ctx context.Context, subscription *repository.Subscription, userID string, subscriptionID int) error {
	args := m.Called(ctx, subscription, userID, subscriptionID)
	return args.Error(0)
}

func (m *MockSubscriptionsRepository) DeleteSubscription(ctx context.Context, userID string, subscriptionID int) error {
	args := m.Called(ctx, userID, subscriptionID)
	return args.Error(0)
}

func (m *MockSubscriptionsRepository) GetSubscriptionsByUserID(ctx context.Context, userID string) ([]*repository.Subscription, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repository.Subscription), args.Error(1)
}

func (m *MockSubscriptionsRepository) GetSubscriptionsByPeriod(ctx context.Context, userID string, serviceNames []string, startDate, endDate time.Time) ([]*repository.Subscription, error) {
	args := m.Called(ctx, userID, serviceNames, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repository.Subscription), args.Error(1)
}

func (m *MockSubscriptionsRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSubscriptionsRepository) RunMigrations(migrationsFilePath string) error {
	args := m.Called(migrationsFilePath)
	return args.Error(0)
}

type SubscriptionServiceTestSuite struct {
	suite.Suite
	mockRepo *MockSubscriptionsRepository
	service  SubscriptionService
}

func (suite *SubscriptionServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockSubscriptionsRepository)
	suite.service = NewSubscriptionService(suite.mockRepo)
}

func (suite *SubscriptionServiceTestSuite) TestCreateSubscription_Success() {
	ctx := context.Background()
	req := &CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      "550e8400-e29b-41d4-a716-446655440000",
		StartDate:   "01-2025",
		EndDate:     "12-2025",
	}

	suite.mockRepo.On("Create", ctx, mock.AnythingOfType("*repository.Subscription")).Return(nil)

	result, err := suite.service.CreateSubscription(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 1, result.ID) // Mock sets ID to 1
	assert.Equal(suite.T(), "Netflix", result.ServiceName)
	assert.Equal(suite.T(), 599, result.Price)
	assert.Equal(suite.T(), req.UserID, result.UserID)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestCreateSubscription_WithoutEndDate() {
	ctx := context.Background()
	req := &CreateSubscriptionRequest{
		ServiceName: "Spotify",
		Price:       299,
		UserID:      "550e8400-e29b-41d4-a716-446655440000",
		StartDate:   "01-2025",
		EndDate:     "",
	}

	suite.mockRepo.On("Create", ctx, mock.AnythingOfType("*repository.Subscription")).Return(nil)

	result, err := suite.service.CreateSubscription(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Nil(suite.T(), result.EndDate)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestCreateSubscription_InvalidUserID() {
	ctx := context.Background()
	req := &CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      "invalid-uuid",
		StartDate:   "01-2025",
	}

	result, err := suite.service.CreateSubscription(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "validation failed")
	suite.mockRepo.AssertNotCalled(suite.T(), "Create")
}

func (suite *SubscriptionServiceTestSuite) TestCreateSubscription_InvalidPrice() {
	ctx := context.Background()
	req := &CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       -100,
		UserID:      "550e8400-e29b-41d4-a716-446655440000",
		StartDate:   "01-2025",
	}

	result, err := suite.service.CreateSubscription(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "validation failed")
	suite.mockRepo.AssertNotCalled(suite.T(), "Create")
}

func (suite *SubscriptionServiceTestSuite) TestCreateSubscription_InvalidDateFormat() {
	ctx := context.Background()
	req := &CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      "550e8400-e29b-41d4-a716-446655440000",
		StartDate:   "invalid-date",
	}

	result, err := suite.service.CreateSubscription(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInvalidDateFormat, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "Create")
}

func (suite *SubscriptionServiceTestSuite) TestCreateSubscription_EndDateBeforeStart() {
	ctx := context.Background()
	req := &CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      "550e8400-e29b-41d4-a716-446655440000",
		StartDate:   "12-2025",
		EndDate:     "01-2025",
	}

	result, err := suite.service.CreateSubscription(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrEndDateBeforeStart, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "Create")
}

func (suite *SubscriptionServiceTestSuite) TestCreateSubscription_RepositoryError() {
	ctx := context.Background()
	req := &CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      "550e8400-e29b-41d4-a716-446655440000",
		StartDate:   "01-2025",
	}

	suite.mockRepo.On("Create", ctx, mock.AnythingOfType("*repository.Subscription")).Return(errors.New("database error"))

	result, err := suite.service.CreateSubscription(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInternalServer, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestGetSubscription_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1

	expectedSub := &repository.Subscription{
		ID:          subscriptionID,
		ServiceName: "Netflix",
		Price:       599,
		UserID:      userID,
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	suite.mockRepo.On("GetSubscription", ctx, userID, subscriptionID).Return(expectedSub, nil)

	result, err := suite.service.GetSubscription(ctx, userID, subscriptionID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedSub, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestGetSubscription_InvalidUserID() {
	ctx := context.Background()
	userID := "invalid-uuid"
	subscriptionID := 1

	result, err := suite.service.GetSubscription(ctx, userID, subscriptionID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInvalidUserID, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "GetSubscription")
}

func (suite *SubscriptionServiceTestSuite) TestGetSubscription_InvalidSubscriptionID() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 0

	result, err := suite.service.GetSubscription(ctx, userID, subscriptionID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInvalidSubscriptionID, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "GetSubscription")
}

func (suite *SubscriptionServiceTestSuite) TestGetSubscription_NotFound() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 999

	suite.mockRepo.On("GetSubscription", ctx, userID, subscriptionID).Return(nil, errors.New("not found"))

	result, err := suite.service.GetSubscription(ctx, userID, subscriptionID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrSubscriptionNotFound, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestUpdateSubscription_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1

	req := &UpdateSubscriptionRequest{
		ServiceName: "Netflix Premium",
		Price:       799,
		StartDate:   "01-2025",
		EndDate:     "12-2025",
	}

	suite.mockRepo.On("UpdateSubscription", ctx, mock.AnythingOfType("*repository.Subscription"), userID, subscriptionID).Return(nil)

	result, err := suite.service.UpdateSubscription(ctx, userID, subscriptionID, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), subscriptionID, result.ID)
	assert.Equal(suite.T(), userID, result.UserID)
	assert.Equal(suite.T(), "Netflix Premium", result.ServiceName)
	assert.Equal(suite.T(), 799, result.Price)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestUpdateSubscription_InvalidUserID() {
	ctx := context.Background()
	userID := "invalid-uuid"
	subscriptionID := 1
	req := &UpdateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		StartDate:   "01-2025",
	}

	result, err := suite.service.UpdateSubscription(ctx, userID, subscriptionID, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInvalidUserID, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "UpdateSubscription")
}

func (suite *SubscriptionServiceTestSuite) TestUpdateSubscription_InvalidSubscriptionID() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := -1
	req := &UpdateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		StartDate:   "01-2025",
	}

	result, err := suite.service.UpdateSubscription(ctx, userID, subscriptionID, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInvalidSubscriptionID, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "UpdateSubscription")
}

func (suite *SubscriptionServiceTestSuite) TestUpdateSubscription_RepositoryError() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1
	req := &UpdateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       599,
		StartDate:   "01-2025",
	}

	suite.mockRepo.On("UpdateSubscription", ctx, mock.AnythingOfType("*repository.Subscription"), userID, subscriptionID).Return(errors.New("update failed"))

	result, err := suite.service.UpdateSubscription(ctx, userID, subscriptionID, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrSubscriptionNotFound, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestDeleteSubscription_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1

	suite.mockRepo.On("DeleteSubscription", ctx, userID, subscriptionID).Return(nil)

	err := suite.service.DeleteSubscription(ctx, userID, subscriptionID)

	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestDeleteSubscription_InvalidUserID() {
	ctx := context.Background()
	userID := "invalid-uuid"
	subscriptionID := 1

	err := suite.service.DeleteSubscription(ctx, userID, subscriptionID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidUserID, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "DeleteSubscription")
}

func (suite *SubscriptionServiceTestSuite) TestDeleteSubscription_InvalidSubscriptionID() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 0

	err := suite.service.DeleteSubscription(ctx, userID, subscriptionID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidSubscriptionID, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "DeleteSubscription")
}

func (suite *SubscriptionServiceTestSuite) TestDeleteSubscription_NotFound() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 999

	suite.mockRepo.On("DeleteSubscription", ctx, userID, subscriptionID).Return(errors.New("not found"))

	err := suite.service.DeleteSubscription(ctx, userID, subscriptionID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrSubscriptionNotFound, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestGetUserSubscriptions_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"

	expectedSubs := []*repository.Subscription{
		{
			ID:          1,
			ServiceName: "Netflix",
			Price:       599,
			UserID:      userID,
			StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			ServiceName: "Spotify",
			Price:       299,
			UserID:      userID,
			StartDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	suite.mockRepo.On("GetSubscriptionsByUserID", ctx, userID).Return(expectedSubs, nil)

	result, err := suite.service.GetUserSubscriptions(ctx, userID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), expectedSubs, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestGetUserSubscriptions_InvalidUserID() {
	ctx := context.Background()
	userID := "invalid-uuid"

	result, err := suite.service.GetUserSubscriptions(ctx, userID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInvalidUserID, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "GetSubscriptionsByUserID")
}

func (suite *SubscriptionServiceTestSuite) TestGetUserSubscriptions_RepositoryError() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"

	suite.mockRepo.On("GetSubscriptionsByUserID", ctx, userID).Return(nil, errors.New("database error"))

	result, err := suite.service.GetUserSubscriptions(ctx, userID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInternalServer, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestCalculateTotalCost_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"

	req := &GetCostRequest{
		UserID:    userID,
		StartDate: "01-2025",
		EndDate:   "06-2025",
	}

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 23, 59, 59, 999999999, time.UTC)

	subscriptions := []*repository.Subscription{
		{
			ID:          1,
			ServiceName: "Netflix",
			Price:       599,
			UserID:      userID,
			StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     nil, // Active subscription
		},
		{
			ID:          2,
			ServiceName: "Spotify",
			Price:       299,
			UserID:      userID,
			StartDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     nil,
		},
	}

	suite.mockRepo.On("GetSubscriptionsByPeriod", ctx, userID, []string(nil), startDate, endDate).Return(subscriptions, nil)

	result, err := suite.service.CalculateTotalCost(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), userID, result.UserID)
	assert.Equal(suite.T(), "01-2025", result.StartDate)
	assert.Equal(suite.T(), "06-2025", result.EndDate)
	assert.Greater(suite.T(), result.TotalCost, 0)
	assert.Len(suite.T(), result.Breakdown, 2)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestCalculateTotalCost_WithServiceFilter() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"

	req := &GetCostRequest{
		UserID:       userID,
		ServiceNames: []string{"Netflix"},
		StartDate:    "01-2025",
		EndDate:      "06-2025",
	}

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 23, 59, 59, 999999999, time.UTC)

	subscriptions := []*repository.Subscription{
		{
			ID:          1,
			ServiceName: "Netflix",
			Price:       599,
			UserID:      userID,
			StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     nil,
		},
	}

	suite.mockRepo.On("GetSubscriptionsByPeriod", ctx, userID, []string{"Netflix"}, startDate, endDate).Return(subscriptions, nil)

	result, err := suite.service.CalculateTotalCost(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Breakdown, 1)
	assert.Equal(suite.T(), "Netflix", result.Breakdown[0].ServiceName)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SubscriptionServiceTestSuite) TestCalculateTotalCost_InvalidUserID() {
	ctx := context.Background()
	req := &GetCostRequest{
		UserID:    "invalid-uuid",
		StartDate: "01-2025",
		EndDate:   "06-2025",
	}

	result, err := suite.service.CalculateTotalCost(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "validation failed")
	suite.mockRepo.AssertNotCalled(suite.T(), "GetSubscriptionsByPeriod")
}

func (suite *SubscriptionServiceTestSuite) TestCalculateTotalCost_InvalidDateFormat() {
	ctx := context.Background()
	req := &GetCostRequest{
		UserID:    "550e8400-e29b-41d4-a716-446655440000",
		StartDate: "invalid-date",
		EndDate:   "06-2025",
	}

	result, err := suite.service.CalculateTotalCost(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInvalidDateFormat, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "GetSubscriptionsByPeriod")
}

func (suite *SubscriptionServiceTestSuite) TestCalculateTotalCost_InvalidDateRange() {
	ctx := context.Background()
	req := &GetCostRequest{
		UserID:    "550e8400-e29b-41d4-a716-446655440000",
		StartDate: "06-2025",
		EndDate:   "01-2025", // End before start
	}

	result, err := suite.service.CalculateTotalCost(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInvalidDateRange, err)
	suite.mockRepo.AssertNotCalled(suite.T(), "GetSubscriptionsByPeriod")
}

func (suite *SubscriptionServiceTestSuite) TestCalculateTotalCost_RepositoryError() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	req := &GetCostRequest{
		UserID:    userID,
		StartDate: "01-2025",
		EndDate:   "06-2025",
	}

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 6, 30, 23, 59, 59, 999999999, time.UTC)

	suite.mockRepo.On("GetSubscriptionsByPeriod", ctx, userID, []string(nil), startDate, endDate).Return(nil, errors.New("database error"))

	result, err := suite.service.CalculateTotalCost(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrInternalServer, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestSubscriptionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionServiceTestSuite))
}
