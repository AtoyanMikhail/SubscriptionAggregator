package handlers

import (
	"time"

	"github.com/AtoyanMikhail/SubscribtionAggregation/internal/logger"
	"github.com/gin-gonic/gin"
)

// LoggerMiddleware creates a custom logger middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log := logger.Global()
		log.Info("HTTP request completed",
			logger.String("method", method),
			logger.String("path", path),
			logger.String("client_ip", clientIP),
			logger.Int("status_code", statusCode),
			logger.Any("latency", latency),
			logger.Int("body_size", c.Writer.Size()),
		)

		if len(c.Errors) > 0 {
			log.Error("HTTP request errors",
				logger.String("path", path),
				logger.Any("errors", c.Errors.String()))
		}
	}
}

// ErrorHandlerMiddleware handles panics and converts them to proper error responses
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Global().Error("panic recovered in HTTP handler",
					logger.Any("error", err),
					logger.String("path", c.Request.URL.Path),
					logger.String("method", c.Request.Method))

				c.JSON(500, ErrorResponse{
					Error:   "internal server error",
					Message: "an unexpected error occurred",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
