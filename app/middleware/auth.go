package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services/auth"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
)

const (
	contextUserIDKey    = "userID"
	contextSessionIDKey = "sessionID"
)

// AuthRequired middleware verifies JWT access token and sets user context
// Supports both new token format (with sub, sid, role_id, role_name) and old format (with user_id)
// This ensures backward compatibility while using the modern auth service
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			c.Abort()
			return
		}

		// Extract Bearer token - handle both "Bearer <token>" and raw token formats
		var tokenString string
		authHeader = strings.TrimSpace(authHeader)
		
		// Check if it starts with "Bearer " (case-insensitive)
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
				c.Abort()
				return
			}
			tokenString = strings.TrimSpace(parts[1])
		} else {
			// If no "Bearer " prefix, assume the entire header is the token
			// This supports Swagger UI's ApiKeyAuth which sends raw tokens
			tokenString = authHeader
		}
		
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			c.Abort()
			return
		}

		// Verify and parse token
		claims, err := auth.VerifyAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Extract user ID - support both new format (sub) and old format (user_id)
		var userID int64
		if sub, ok := claims["sub"].(string); ok {
			// New token format - sub contains user ID as string
			_, err := fmt.Sscanf(sub, "%d", &userID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
				c.Abort()
				return
			}
		} else if userIDFloat, ok := claims["user_id"].(float64); ok {
			// Old token format - user_id as float64
			userID = int64(userIDFloat)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims: missing user_id"})
			c.Abort()
			return
		}

		// Extract session ID (optional for backward compatibility)
		sessionID := ""
		if sid, ok := claims["sid"].(string); ok {
			sessionID = sid
		}

		// Extract role information (optional for backward compatibility)
		var roleID int64
		var roleName string
		
		// Try to get role_id from token
		if roleIDFloat, ok := claims["role_id"].(float64); ok {
			roleID = int64(roleIDFloat)
		} else if roleIDInt, ok := claims["role_id"].(int64); ok {
			roleID = roleIDInt
		}
		
		// Try to get role_name from token
		if rn, ok := claims["role_name"].(string); ok {
			roleName = rn
		}

		// If role info not in token, try to load from database (backward compatibility)
		if roleID == 0 || roleName == "" {
			// This is handled by the old AuthMiddleware approach
			// For now, we'll allow it but log a warning
			// In production, you may want to require role info in token
		}

		// Set context values
		c.Set(contextUserIDKey, userID)
		if sessionID != "" {
			c.Set(contextSessionIDKey, sessionID)
		}
		if roleID > 0 {
			c.Set("roleID", roleID)
		}
		if roleName != "" {
			c.Set("roleName", roleName)
		}

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

// GetUserEmail extracts user email from gin context
// This function queries the database to get the user's email based on the user ID
func GetUserEmail(c *gin.Context) (string, bool) {
	userID, exists := GetUserID(c)
	if !exists {
		return "", false
	}
	
	// Query database to get user email
	var email string
	err := config.DB.Model(&models.User{}).
		Where("id = ? AND is_deleted = false", userID).
		Select("email").
		Scan(&email).Error
	
	if err != nil {
		return "", false
	}
	
	return email, true
}




