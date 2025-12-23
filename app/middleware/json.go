package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// StrictJSONBinding is a convenience wrapper for StrictJSONDecoder with a default 1MB limit
func StrictJSONBinding() gin.HandlerFunc {
	return StrictJSONDecoder(1 << 20) // 1 MB default limit
}

// StrictJSONDecoder ensures that JSON decoding disallows unknown fields
// and sets a size limit. It reads and validates the body, then restores it for handlers.
func StrictJSONDecoder(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")
		if contentType != "application/json" && !strings.HasPrefix(contentType, "application/json;") {
			c.Next()
			return
		}

		// Limit request body size
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		// Read the body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
			c.Abort()
			return
		}

		// Restore the body for handlers to read
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Validate JSON and check for unknown fields
		if len(bodyBytes) > 0 {
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			decoder.DisallowUnknownFields()

			var body map[string]interface{}
			if err := decoder.Decode(&body); err != nil {
				if err != io.EOF {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json or unknown fields", "details": err.Error()})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

