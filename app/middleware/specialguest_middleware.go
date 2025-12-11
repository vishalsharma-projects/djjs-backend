package middleware

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
)

// ValidateSpecialGuestMiddleware checks if the special guest exists and user has permission
func ValidateSpecialGuestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip validation for GET all and POST routes
		if c.Request.Method == http.MethodGet && c.Param("id") == "" {
			c.Next()
			return
		}

		if c.Request.Method == http.MethodPost {
			c.Next()
			return
		}

		// Get and validate special guest ID
		specialGuestID := c.Param("id")
		if specialGuestID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "special guest ID is required"})
			c.Abort()
			return
		}

		sgID, err := strconv.ParseUint(specialGuestID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid special guest ID format"})
			c.Abort()
			return
		}

		// Check if special guest exists
		var specialGuest models.SpecialGuest
		if err := config.DB.First(&specialGuest, sgID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "special guest not found",
				"id":    specialGuestID,
			})
			c.Abort()
			return
		}

		// Validate user authentication
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		roleID, exists := c.Get("roleID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user role not found"})
			c.Abort()
			return
		}

		// Role-based access control
		role, ok := roleID.(uint)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid role type"})
			c.Abort()
			return
		}

		switch role {
		case 1, 2: // Admin, Manager - full access
			// Allow all operations
		case 3: // Regular user - read only
			if c.Request.Method != http.MethodGet {
				c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
				c.Abort()
				return
			}
		default:
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid role"})
			c.Abort()
			return
		}

		// Set context values for handlers
		c.Set("specialGuest", &specialGuest)
		c.Set("currentUserID", userID)
		c.Next()
	}
}
