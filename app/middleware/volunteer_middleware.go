package middleware

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ValidateVolunteerMiddleware loads volunteer by ID and shares it via context
func ValidateVolunteerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		volunteerIDParam := c.Param("id")
		if volunteerIDParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "volunteer id is required"})
			c.Abort()
			return
		}

		volunteerID, err := strconv.ParseUint(volunteerIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid volunteer id format"})
			c.Abort()
			return
		}

		var volunteer models.Volunteer
		if err := config.DB.Preload("Branch").First(&volunteer, uint(volunteerID)).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "volunteer not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch volunteer"})
			}
			c.Abort()
			return
		}

		c.Set("volunteer", &volunteer)
		c.Next()
	}
}
