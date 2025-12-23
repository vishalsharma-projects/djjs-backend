package middleware

import (
	"net/http"
	"strings"

	"github.com/followCode/djjs-event-reporting-backend/app/services/auth"
	"github.com/gin-gonic/gin"
)

const (
	contextUserIDKey    = "userID"
	contextSessionIDKey = "sessionID"
)

// AuthRequired middleware verifies JWT access token and sets user context
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verify and parse token
		claims, err := auth.VerifyAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Extract user ID
		userID, err := auth.ParseUserIDFromToken(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Extract session ID
		sessionID, err := auth.ParseSessionIDFromToken(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Set context values
		c.Set(contextUserIDKey, userID)
		c.Set(contextSessionIDKey, sessionID)

		c.Next()
	}
}

// GetUserID extracts user ID from gin context
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(contextUserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

// GetSessionID extracts session ID from gin context
func GetSessionID(c *gin.Context) (string, bool) {
	sessionID, exists := c.Get(contextSessionIDKey)
	if !exists {
		return "", false
	}
	sid, ok := sessionID.(string)
	return sid, ok
}


