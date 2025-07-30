package service

import (
	"time"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/repository"
)

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" validate:"required,min=1,max=255"`
	Price       int    `json:"price" validate:"required,min=0"`
	UserID      string `json:"user_id" validate:"required,uuid4"`
	StartDate   string `json:"start_date" validate:"required"` // Format: MM-YYYY
	EndDate     string `json:"end_date,omitempty"`             // Format: MM-YYYY, optional
}

type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name" validate:"required,min=1,max=255"`
	Price       int    `json:"price" validate:"required,min=0"`
	StartDate   string `json:"start_date" validate:"required"` // Format: MM-YYYY
	EndDate     string `json:"end_date,omitempty"`             // Format: MM-YYYY, optional
}

type GetCostRequest struct {
	UserID       string   `json:"user_id" validate:"required,uuid4"`
	ServiceNames []string `json:"service_names,omitempty"`        // Optional filter
	StartDate    string   `json:"start_date" validate:"required"` // Format: MM-YYYY
	EndDate      string   `json:"end_date" validate:"required"`   // Format: MM-YYYY
}

type CostResponse struct {
	UserID    string                      `json:"user_id"`
	StartDate string                      `json:"start_date"`
	EndDate   string                      `json:"end_date"`
	TotalCost int                         `json:"total_cost"`
	Breakdown []SubscriptionCostBreakdown `json:"breakdown"`
}

type SubscriptionCostBreakdown struct {
	SubscriptionID int    `json:"subscription_id"`
	ServiceName    string `json:"service_name"`
	MonthlyPrice   int    `json:"monthly_price"`
	MonthsCount    int    `json:"months_count"`
	TotalCost      int    `json:"total_cost"`
}

func (r *CreateSubscriptionRequest) ToSubscriptionModel() (*repository.Subscription, error) {
	startDate, err := ParseMonthYear(r.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if r.EndDate != "" {
		endDateTime, err := ParseMonthYear(r.EndDate)
		if err != nil {
			return nil, err
		}
		// Set end date to last day of the month
		lastDay := GetLastDayOfMonth(endDateTime)
		endDate = &lastDay
	}

	// Validate that end date is after start date
	if endDate != nil && endDate.Before(startDate) {
		return nil, ErrEndDateBeforeStart
	}

	return &repository.Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      r.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

// ToSubscriptionModel converts UpdateSubscriptionRequest to Subscription model
func (r *UpdateSubscriptionRequest) ToSubscriptionModel() (*repository.Subscription, error) {
	startDate, err := ParseMonthYear(r.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if r.EndDate != "" {
		endDateTime, err := ParseMonthYear(r.EndDate)
		if err != nil {
			return nil, err
		}
		lastDay := GetLastDayOfMonth(endDateTime)
		endDate = &lastDay
	}

	if endDate != nil && endDate.Before(startDate) {
		return nil, ErrEndDateBeforeStart
	}

	return &repository.Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
