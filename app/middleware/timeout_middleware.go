package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware creates a middleware that sets a timeout for requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace the request context with the timeout context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}



