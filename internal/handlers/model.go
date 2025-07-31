package handlers

import (
	"time"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/repository"
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/service"
)

// HTTP Request/Response models for Swagger documentation

// CreateSubscriptionRequest represents the request body for creating a subscription
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int    `json:"price" binding:"required,min=0" example:"400"`
	UserID      string `json:"user_id" binding:"required,uuid4" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date" binding:"required" example:"07-2025"`
	EndDate     string `json:"end_date,omitempty" example:"12-2025"`
} // @name CreateSubscriptionRequest

// UpdateSubscriptionRequest represents the request body for updating a subscription
type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required" example:"Netflix Premium"`
	Price       int    `json:"price" binding:"required,min=0" example:"599"`
	StartDate   string `json:"start_date" binding:"required" example:"07-2025"`
	EndDate     string `json:"end_date,omitempty" example:"12-2025"`
} // @name UpdateSubscriptionRequest

// GetCostRequest represents the request body/query params for calculating total cost
type GetCostRequest struct {
	UserID       string   `json:"user_id" form:"user_id" binding:"required,uuid4" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceNames []string `json:"service_names,omitempty" form:"service_names" example:"Netflix,Spotify"`
	StartDate    string   `json:"start_date" form:"start_date" binding:"required" example:"01-2025"`
	EndDate      string   `json:"end_date" form:"end_date" binding:"required" example:"12-2025"`
} // @name GetCostRequest

// SubscriptionResponse represents a subscription in API responses
type SubscriptionResponse struct {
	ID          int     `json:"id" example:"1"`
	ServiceName string  `json:"service_name" example:"Yandex Plus"`
	Price       int     `json:"price" example:"400"`
	UserID      string  `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string  `json:"start_date" example:"2025-07-01T00:00:00Z"`
	EndDate     *string `json:"end_date,omitempty" example:"2025-12-31T23:59:59Z"`
} // @name SubscriptionResponse

// CostResponse represents the response for cost calculation
type CostResponse struct {
	UserID    string                      `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate string                      `json:"start_date" example:"01-2025"`
	EndDate   string                      `json:"end_date" example:"12-2025"`
	TotalCost int                         `json:"total_cost" example:"4800"`
	Breakdown []SubscriptionCostBreakdown `json:"breakdown"`
} // @name CostResponse

// SubscriptionCostBreakdown represents cost breakdown for each subscription
type SubscriptionCostBreakdown struct {
	SubscriptionID int    `json:"subscription_id" example:"1"`
	ServiceName    string `json:"service_name" example:"Netflix"`
	MonthlyPrice   int    `json:"monthly_price" example:"599"`
	MonthsCount    int    `json:"months_count" example:"6"`
	TotalCost      int    `json:"total_cost" example:"3594"`
} // @name SubscriptionCostBreakdown

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"validation failed"`
	Message string `json:"message,omitempty" example:"invalid user ID format"`
} // @name ErrorResponse

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string `json:"message" example:"operation completed successfully"`
} // @name SuccessResponse

// ListSubscriptionsResponse represents response for listing subscriptions
type ListSubscriptionsResponse struct {
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
	Count         int                    `json:"count" example:"5"`
} // @name ListSubscriptionsResponse

// Convert service request to handler request
func (r *CreateSubscriptionRequest) ToServiceRequest() *service.CreateSubscriptionRequest {
	return &service.CreateSubscriptionRequest{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      r.UserID,
		StartDate:   r.StartDate,
		EndDate:     r.EndDate,
	}
}

func (r *UpdateSubscriptionRequest) ToServiceRequest() *service.UpdateSubscriptionRequest {
	return &service.UpdateSubscriptionRequest{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		StartDate:   r.StartDate,
		EndDate:     r.EndDate,
	}
}

func (r *GetCostRequest) ToServiceRequest() *service.GetCostRequest {
	return &service.GetCostRequest{
		UserID:       r.UserID,
		ServiceNames: r.ServiceNames,
		StartDate:    r.StartDate,
		EndDate:      r.EndDate,
	}
}

// Convert model to response
func SubscriptionToResponse(sub *repository.Subscription) SubscriptionResponse {
	resp := SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format(time.RFC3339),
	}

	if sub.EndDate != nil {
		endDateStr := sub.EndDate.Format(time.RFC3339)
		resp.EndDate = &endDateStr
	}

	return resp
}

func SubscriptionsToResponse(subs []*repository.Subscription) ListSubscriptionsResponse {
	responses := make([]SubscriptionResponse, len(subs))
	for i, sub := range subs {
		responses[i] = SubscriptionToResponse(sub)
	}

	return ListSubscriptionsResponse{
		Subscriptions: responses,
		Count:         len(responses),
	}
}

func ServiceCostToResponse(serviceCost *service.CostResponse) CostResponse {
	breakdown := make([]SubscriptionCostBreakdown, len(serviceCost.Breakdown))
	for i, item := range serviceCost.Breakdown {
		breakdown[i] = SubscriptionCostBreakdown{
			SubscriptionID: item.SubscriptionID,
			ServiceName:    item.ServiceName,
			MonthlyPrice:   item.MonthlyPrice,
			MonthsCount:    item.MonthsCount,
			TotalCost:      item.TotalCost,
		}
	}

	return CostResponse{
		UserID:    serviceCost.UserID,
		StartDate: serviceCost.StartDate,
		EndDate:   serviceCost.EndDate,
		TotalCost: serviceCost.TotalCost,
		Breakdown: breakdown,
	}
}
