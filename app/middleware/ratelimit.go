package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig configures rate limiting behavior
type RateLimitConfig struct {
	MaxRequests     int
	Window          time.Duration
	IdentifierKey   string // "ip" or a custom key from context
	CustomIdentifier func(*gin.Context) string // Custom function to extract identifier
}

// RateLimiter creates a rate limiting middleware
func RateLimiter(cfg RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.RedisClient == nil {
			// If Redis is not available, skip rate limiting
			c.Next()
			return
		}

		ctx := context.Background()

		// Get identifier
		var identifier string
		if cfg.CustomIdentifier != nil {
			identifier = cfg.CustomIdentifier(c)
		} else if cfg.IdentifierKey == "ip" {
			identifier = GetClientIP(c)
		} else {
			identifier = c.GetString(cfg.IdentifierKey)
		}

		if identifier == "" {
			c.Next()
			return
		}

		// Create Redis key
		key := fmt.Sprintf("ratelimit:%s:%s", c.FullPath(), identifier)

		// Check current count
		count, err := config.RedisClient.Get(ctx, key).Int()
		if err == redis.Nil {
			count = 0
		} else if err != nil {
			// Redis error - allow request but log
			c.Next()
			return
		}

		// Check if limit exceeded
		if count >= cfg.MaxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		// Increment counter
		pipe := config.RedisClient.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, cfg.Window)
		_, err = pipe.Exec(ctx)
		if err != nil {
			// Redis error - allow request
			c.Next()
			return
		}

		c.Next()
	}
}

// GetClientIP extracts the client IP, respecting X-Forwarded-For if trusted proxy
func GetClientIP(c *gin.Context) string {
	if config.TrustProxy {
		// Check X-Forwarded-For header
		if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
			// X-Forwarded-For can contain multiple IPs, take the first one
			for i := 0; i < len(forwarded); i++ {
				if forwarded[i] == ',' {
					return forwarded[:i]
				}
			}
			return forwarded
		}
		// Check X-Real-IP
		if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
			return realIP
		}
	}
	return c.ClientIP()
}


