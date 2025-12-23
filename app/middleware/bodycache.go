package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

// BodyCache caches the request body so it can be read multiple times
func BodyCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			// Read body
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.Next()
				return
			}

			// Restore body for subsequent reads
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// Store in context
			c.Set("cached_body", body)
		}
		c.Next()
	}
}

// GetCachedBody retrieves cached body from context
func GetCachedBody(c *gin.Context) []byte {
	body, exists := c.Get("cached_body")
	if !exists {
		return nil
	}
	b, ok := body.([]byte)
	if !ok {
		return nil
	}
	return b
}


