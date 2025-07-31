package handlers

import (
	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(subscriptionService service.SubscriptionService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Health check endpoint
	router.GET("/health", HealthCheck)

	subscriptionHandler := NewSubscriptionHandler(subscriptionService)

	v1 := router.Group("/api/v1")
	{
		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.CreateSubscription)
			subscriptions.GET("/:user_id/:subscription_id", subscriptionHandler.GetSubscription)
			subscriptions.PUT("/:user_id/:subscription_id", subscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:user_id/:subscription_id", subscriptionHandler.DeleteSubscription)
			subscriptions.GET("/user/:user_id", subscriptionHandler.GetUserSubscriptions)
			subscriptions.GET("/cost", subscriptionHandler.CalculateTotalCostQuery)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
