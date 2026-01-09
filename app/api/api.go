package api

import (
	"context"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes and groups them together
func SetupRoutes(r *gin.Engine) {
	// Health check endpoint (public, no auth required)
	r.GET("/health", HealthCheckHandler)
	r.GET("/api/health", HealthCheckHandler)

	// Main API group
	api := r.Group("/api")
	{
		// Authentication routes
		SetupAuthRoutes(api)

		// CRUD routes
		SetupAreaRoutes(api)
		SetupUserRoutes(api)
		SetupBranchRoutes(api)
		SetupChildBranchRoutes(api)
		SetupEventRoutes(api)
		SetupPromotionRoutes(api)
		SetupMediaRoutes(api)
		SetupSpecialGuestRoutes(api)
		SetupVolunteerRoutes(api)
		SetupDonationRoutes(api)
		SetupMasterRoutes(api)
		SetupFileRoutes(api)
		SetupBranchMediaRoutes(api)
		SetupChildBranchMediaRoutes(api)
		SetupUserPreferencesRoutes(api)

		// RBAC routes
		rbacHandler := handlers.NewRBACHandler(config.DB)
		SetupRBACRoutes(api, rbacHandler)
	}
}

// HealthCheckHandler returns the health status of the API including S3 connectivity
// @Summary Health check endpoint
// @Description Returns the health status of the API and S3 bucket connectivity
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "Health status"
// @Router /health [get]
// @Router /api/health [get]
func HealthCheckHandler(c *gin.Context) {
	health := gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"services":  make(map[string]interface{}),
	}

	servicesMap := health["services"].(map[string]interface{})

	// Check S3 connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := services.VerifyS3Connection(ctx); err != nil {
		servicesMap["s3"] = gin.H{
			"status":  "error",
			"message": err.Error(),
		}
		health["status"] = "degraded"
	} else {
		servicesMap["s3"] = gin.H{
			"status":  "ok",
			"message": "S3 bucket is accessible and has correct permissions",
		}
	}

	statusCode := 200
	if health["status"] == "degraded" {
		statusCode = 503 // Service Unavailable
	}

	c.JSON(statusCode, health)
}

