package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/logger"
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(subscriptionService service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// CreateSubscription creates a new subscription
// @Summary Create a new subscription
// @Description Create a new subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Global().Error("failed to bind create subscription request", logger.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation failed",
			Message: err.Error(),
		})
		return
	}

	subscription, err := h.subscriptionService.CreateSubscription(c.Request.Context(), req.ToServiceRequest())
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := SubscriptionToResponse(subscription)
	c.JSON(http.StatusCreated, response)
}

// GetSubscription retrieves a specific subscription
// @Summary Get a subscription by ID
// @Description Get a specific subscription for a user
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param subscription_id path int true "Subscription ID"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/subscriptions/{user_id}/{subscription_id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	userID := c.Param("user_id")
	subscriptionIDStr := c.Param("subscription_id")

	subscriptionID, err := strconv.Atoi(subscriptionIDStr)
	if err != nil {
		logger.Global().Error("invalid subscription ID", logger.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid subscription ID",
			Message: "subscription ID must be a valid integer",
		})
		return
	}

	subscription, err := h.subscriptionService.GetSubscription(c.Request.Context(), userID, subscriptionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := SubscriptionToResponse(subscription)
	c.JSON(http.StatusOK, response)
}

// UpdateSubscription updates an existing subscription
// @Summary Update a subscription
// @Description Update an existing subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param subscription_id path int true "Subscription ID"
// @Param subscription body UpdateSubscriptionRequest true "Updated subscription data"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/subscriptions/{user_id}/{subscription_id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	userID := c.Param("user_id")
	subscriptionIDStr := c.Param("subscription_id")

	subscriptionID, err := strconv.Atoi(subscriptionIDStr)
	if err != nil {
		logger.Global().Error("invalid subscription ID", logger.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid subscription ID",
			Message: "subscription ID must be a valid integer",
		})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Global().Error("failed to bind update subscription request", logger.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation failed",
			Message: err.Error(),
		})
		return
	}

	subscription, err := h.subscriptionService.UpdateSubscription(c.Request.Context(), userID, subscriptionID, req.ToServiceRequest())
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := SubscriptionToResponse(subscription)
	c.JSON(http.StatusOK, response)
}

// DeleteSubscription deletes a subscription
// @Summary Delete a subscription
// @Description Delete a specific subscription for a user
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param subscription_id path int true "Subscription ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/subscriptions/{user_id}/{subscription_id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	userID := c.Param("user_id")
	subscriptionIDStr := c.Param("subscription_id")

	subscriptionID, err := strconv.Atoi(subscriptionIDStr)
	if err != nil {
		logger.Global().Error("invalid subscription ID", logger.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid subscription ID",
			Message: "subscription ID must be a valid integer",
		})
		return
	}

	err = h.subscriptionService.DeleteSubscription(c.Request.Context(), userID, subscriptionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "subscription deleted successfully",
	})
}

// GetUserSubscriptions retrieves all subscriptions for a user
// @Summary Get all user subscriptions
// @Description Get all subscriptions for a specific user
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Success 200 {object} ListSubscriptionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/subscriptions/user/{user_id} [get]
func (h *SubscriptionHandler) GetUserSubscriptions(c *gin.Context) {
	userID := c.Param("user_id")

	subscriptions, err := h.subscriptionService.GetUserSubscriptions(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := SubscriptionsToResponse(subscriptions)
	c.JSON(http.StatusOK, response)
}

// CalculateTotalCostQuery calculates total cost using query parameters (alternative endpoint)
// @Summary Calculate total subscription cost (query params)
// @Description Calculate total cost of chosen subscriptions for a user within a specified period using query parameters
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID" format(uuid)
// @Param start_date query string true "Start date in MM-YYYY format"
// @Param end_date query string true "End date in MM-YYYY format"
// @Param service_names query []string false "Service names to filter (comma-separated)"
// @Success 200 {object} CostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/subscriptions/cost [get]
func (h *SubscriptionHandler) CalculateTotalCostQuery(c *gin.Context) {
	var req GetCostRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Global().Error("failed to bind cost calculation query", logger.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation failed",
			Message: err.Error(),
		})
		return
	}

	// Handle comma-separated service names
	serviceNamesParam := c.Query("service_names")
	if serviceNamesParam != "" {
		req.ServiceNames = strings.Split(serviceNamesParam, ",")
		// Trim whitespace from each service name
		for i, name := range req.ServiceNames {
			req.ServiceNames[i] = strings.TrimSpace(name)
		}
	}

	costResponse, err := h.subscriptionService.CalculateTotalCost(c.Request.Context(), req.ToServiceRequest())
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := ServiceCostToResponse(costResponse)
	c.JSON(http.StatusOK, response)
}

// HealthCheck provides a health check endpoint
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} SuccessResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"message": "subscription service is running",
	})
}

// handleError handles service errors and maps them to appropriate HTTP responses
func (h *SubscriptionHandler) handleError(c *gin.Context, err error) {
	logger.Global().Error("handler error", logger.Error(err))

	switch {
	case errors.Is(err, service.ErrSubscriptionNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "subscription not found",
			Message: "the requested subscription does not exist",
		})
	case errors.Is(err, service.ErrInvalidUserID):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid user ID",
			Message: "user ID must be a valid UUID",
		})
	case errors.Is(err, service.ErrInvalidSubscriptionID):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid subscription ID",
			Message: "subscription ID must be a positive integer",
		})
	case errors.Is(err, service.ErrInvalidServiceName):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid service name",
			Message: "service name cannot be empty",
		})
	case errors.Is(err, service.ErrInvalidPrice):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid price",
			Message: "price must be greater than or equal to zero",
		})
	case errors.Is(err, service.ErrInvalidDateFormat):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid date format",
			Message: "date must be in MM-YYYY format",
		})
	case errors.Is(err, service.ErrEndDateBeforeStart):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid date range",
			Message: "end date must be after start date",
		})
	case errors.Is(err, service.ErrInvalidDateRange):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid date range",
			Message: "end date must be after start date",
		})
	case strings.Contains(err.Error(), "validation failed"):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation failed",
			Message: err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal server error",
			Message: "an unexpected error occurred",
		})
	}
}
