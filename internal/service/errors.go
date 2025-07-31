package service

import "errors"

var (
	// Subscription related errors
	ErrSubscriptionNotFound   = errors.New("subscription not found")
	ErrSubscriptionExists     = errors.New("subscription already exists")
	ErrInvalidSubscriptionID  = errors.New("invalid subscription ID")
	
	// Validation errors
	ErrInvalidUserID         = errors.New("invalid user ID format")
	ErrInvalidServiceName    = errors.New("service name cannot be empty")
	ErrInvalidPrice          = errors.New("price must be greater than or equal to zero")
	ErrInvalidDateFormat     = errors.New("invalid date format, expected MM-YYYY")
	ErrEndDateBeforeStart    = errors.New("end date must be after start date")
	ErrInvalidDateRange      = errors.New("invalid date range")
	
	// General errors
	ErrEmptyResult          = errors.New("no subscriptions found")
	ErrInternalServer       = errors.New("internal server error")
)