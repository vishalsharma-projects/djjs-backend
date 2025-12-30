package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupDonationRoutes configures donation CRUD routes
func SetupDonationRoutes(r *gin.RouterGroup) {
	donations := r.Group("/donations")
	donations.Use(middleware.AuthRequired())
	{
		donations.POST("", handlers.CreateDonation)
		donations.GET("", handlers.GetAllDonations)
		donations.PUT("/:id", handlers.UpdateDonation)
		donations.DELETE("/:id", handlers.DeleteDonation)
	}
}

