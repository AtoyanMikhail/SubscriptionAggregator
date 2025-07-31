package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PostgresRepositoryTestSuite struct {
	suite.Suite
	db   *sqlx.DB
	mock sqlmock.Sqlmock
	repo repository.SubscriptionsRepository
}

func (suite *PostgresRepositoryTestSuite) SetupTest() {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(suite.T(), err)
	
	suite.db = sqlx.NewDb(mockDB, "postgres")
	suite.mock = mock
	suite.repo = NewSubscriptionsRepository(suite.db)
}

func (suite *PostgresRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *PostgresRepositoryTestSuite) TestCreate_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	
	subscription := &repository.Subscription{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     &endDate,
	}
	
	expectedQuery := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate, subscription.EndDate).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	
	err := suite.repo.Create(ctx, subscription)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, subscription.ID)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestCreate_WithoutEndDate() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440001"
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	
	subscription := &repository.Subscription{
		ServiceName: "Spotify",
		Price:       299,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     nil,
	}
	
	expectedQuery := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	
	err := suite.repo.Create(ctx, subscription)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, subscription.ID)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestCreate_DatabaseError() {
	ctx := context.Background()
	subscription := &repository.Subscription{
		ServiceName: "Netflix",
		Price:       599,
		UserID:      "550e8400-e29b-41d4-a716-446655440000",
		StartDate:   time.Now(),
	}
	
	expectedQuery := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate, subscription.EndDate).
		WillReturnError(sql.ErrConnDone)
	
	err := suite.repo.Create(ctx, subscription)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrCreateSubscriptionFailed, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestGetSubscription_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	
	expectedQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1 AND id = $2`
	
	rows := sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
		AddRow(subscriptionID, "Netflix", 599, userID, startDate, endDate)
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(userID, subscriptionID).
		WillReturnRows(rows)
	
	result, err := suite.repo.GetSubscription(ctx, userID, subscriptionID)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), subscriptionID, result.ID)
	assert.Equal(suite.T(), "Netflix", result.ServiceName)
	assert.Equal(suite.T(), 599, result.Price)
	assert.Equal(suite.T(), userID, result.UserID)
	assert.Equal(suite.T(), startDate, result.StartDate)
	assert.Equal(suite.T(), endDate, *result.EndDate)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestGetSubscription_NotFound() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 999
	
	expectedQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1 AND id = $2`
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(userID, subscriptionID).
		WillReturnError(sql.ErrNoRows)
	
	result, err := suite.repo.GetSubscription(ctx, userID, subscriptionID)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrSubscriptionNotFound, err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestGetSubscription_DatabaseError() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1
	
	expectedQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1 AND id = $2`
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(userID, subscriptionID).
		WillReturnError(sql.ErrConnDone)
	
	result, err := suite.repo.GetSubscription(ctx, userID, subscriptionID)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrGetSubscriptionFailed, err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestUpdateSubscription_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1
	startDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	
	subscription := &repository.Subscription{
		ID:          subscriptionID,
		ServiceName: "Updated Service",
		Price:       799,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     nil,
	}
	
	expectedQuery := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4
		WHERE user_id = $5 AND id = $6`
	
	suite.mock.ExpectExec(expectedQuery).
		WithArgs(subscription.ServiceName, subscription.Price, subscription.StartDate, subscription.EndDate, userID, subscriptionID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	
	err := suite.repo.UpdateSubscription(ctx, subscription, userID, subscriptionID)
	
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestUpdateSubscription_NotFound() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 999
	
	subscription := &repository.Subscription{
		ServiceName: "Updated Service",
		Price:       799,
		StartDate:   time.Now(),
	}
	
	expectedQuery := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4
		WHERE user_id = $5 AND id = $6`
	
	suite.mock.ExpectExec(expectedQuery).
		WithArgs(subscription.ServiceName, subscription.Price, subscription.StartDate, subscription.EndDate, userID, subscriptionID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	
	err := suite.repo.UpdateSubscription(ctx, subscription, userID, subscriptionID)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrSubscriptionNotFoundForUpdate, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestUpdateSubscription_DatabaseError() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1
	
	subscription := &repository.Subscription{
		ServiceName: "Updated Service",
		Price:       799,
		StartDate:   time.Now(),
	}
	
	expectedQuery := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4
		WHERE user_id = $5 AND id = $6`
	
	suite.mock.ExpectExec(expectedQuery).
		WithArgs(subscription.ServiceName, subscription.Price, subscription.StartDate, subscription.EndDate, userID, subscriptionID).
		WillReturnError(sql.ErrConnDone)
	
	err := suite.repo.UpdateSubscription(ctx, subscription, userID, subscriptionID)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrUpdateSubscriptionFailed, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestDeleteSubscription_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1
	
	expectedQuery := `DELETE FROM subscriptions WHERE user_id = $1 AND id = $2`
	
	suite.mock.ExpectExec(expectedQuery).
		WithArgs(userID, subscriptionID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	
	err := suite.repo.DeleteSubscription(ctx, userID, subscriptionID)
	
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestDeleteSubscription_NotFound() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 999
	
	expectedQuery := `DELETE FROM subscriptions WHERE user_id = $1 AND id = $2`
	
	suite.mock.ExpectExec(expectedQuery).
		WithArgs(userID, subscriptionID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	
	err := suite.repo.DeleteSubscription(ctx, userID, subscriptionID)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrSubscriptionNotFoundForDeletion, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestDeleteSubscription_DatabaseError() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	subscriptionID := 1
	
	expectedQuery := `DELETE FROM subscriptions WHERE user_id = $1 AND id = $2`
	
	suite.mock.ExpectExec(expectedQuery).
		WithArgs(userID, subscriptionID).
		WillReturnError(sql.ErrConnDone)
	
	err := suite.repo.DeleteSubscription(ctx, userID, subscriptionID)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrDeleteSubscriptionFailed, err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestGetSubscriptionsByUserID_Success() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	startDate1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	startDate2 := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	endDate1 := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	
	expectedQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY start_date DESC`
	
	rows := sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
		AddRow(2, "Spotify", 299, userID, startDate2, nil).
		AddRow(1, "Netflix", 599, userID, startDate1, endDate1)
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(userID).
		WillReturnRows(rows)
	
	result, err := suite.repo.GetSubscriptionsByUserID(ctx, userID)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	
	// Check first subscription (should be ordered by start_date DESC)
	assert.Equal(suite.T(), 2, result[0].ID)
	assert.Equal(suite.T(), "Spotify", result[0].ServiceName)
	assert.Equal(suite.T(), 299, result[0].Price)
	assert.Nil(suite.T(), result[0].EndDate)
	
	// Check second subscription
	assert.Equal(suite.T(), 1, result[1].ID)
	assert.Equal(suite.T(), "Netflix", result[1].ServiceName)
	assert.Equal(suite.T(), 599, result[1].Price)
	assert.NotNil(suite.T(), result[1].EndDate)
	assert.Equal(suite.T(), endDate1, *result[1].EndDate)
	
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestGetSubscriptionsByUserID_EmptyResult() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	
	expectedQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY start_date DESC`
	
	rows := sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"})
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(userID).
		WillReturnRows(rows)
	
	result, err := suite.repo.GetSubscriptionsByUserID(ctx, userID)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestGetSubscriptionsByUserID_DatabaseError() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	
	expectedQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY start_date DESC`
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(userID).
		WillReturnError(sql.ErrConnDone)
	
	result, err := suite.repo.GetSubscriptionsByUserID(ctx, userID)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrGetSubscriptionsByUserIDFailed, err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestGetSubscriptionsByPeriod_WithoutServiceFilter() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	
	expectedQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
		AND start_date <= $3
		AND (end_date IS NULL OR end_date >= $2)
		ORDER BY start_date DESC`
	
	subStartDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	subEndDate := time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC)
	
	rows := sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
		AddRow(1, "Netflix", 599, userID, subStartDate, subEndDate)
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(userID, startDate, endDate).
		WillReturnRows(rows)
	
	result, err := suite.repo.GetSubscriptionsByPeriod(ctx, userID, nil, startDate, endDate)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "Netflix", result[0].ServiceName)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestGetSubscriptionsByPeriod_WithServiceFilter() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	serviceNames := []string{"Netflix", "Spotify"}
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	
	expectedQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
		AND start_date <= $3
		AND (end_date IS NULL OR end_date >= $2)
		AND service_name IN ($4,$5)
		ORDER BY start_date DESC`
	
	subStartDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	
	rows := sqlmock.NewRows([]string{"id", "service_name", "price", "user_id", "start_date", "end_date"}).
		AddRow(1, "Netflix", 599, userID, subStartDate, nil).
		AddRow(2, "Spotify", 299, userID, subStartDate, nil)
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(userID, startDate, endDate, "Netflix", "Spotify").
		WillReturnRows(rows)
	
	result, err := suite.repo.GetSubscriptionsByPeriod(ctx, userID, serviceNames, startDate, endDate)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresRepositoryTestSuite) TestGetSubscriptionsByPeriod_DatabaseError() {
	ctx := context.Background()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	
	expectedQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
		AND start_date <= $3
		AND (end_date IS NULL OR end_date >= $2)
		ORDER BY start_date DESC`
	
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(userID, startDate, endDate).
		WillReturnError(sql.ErrConnDone)
	
	result, err := suite.repo.GetSubscriptionsByPeriod(ctx, userID, nil, startDate, endDate)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrGetSubscriptionsByPeriodFailed, err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// Custom matcher for time.Time arguments in mocks
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestPostgresRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepositoryTestSuite))
}