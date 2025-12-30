package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupVolunteerRoutes configures volunteer routes
func SetupVolunteerRoutes(r *gin.RouterGroup) {
	volunteers := r.Group("/volunteers")
	volunteers.Use(middleware.AuthRequired())
	{
		volunteers.POST("", handlers.CreateVolunteerHandler)
		volunteers.GET("", handlers.GetAllVolunteersHandler)
		volunteers.GET("/search", handlers.SearchVolunteersHandler)
		volunteers.PUT("/:id", middleware.ValidateVolunteerMiddleware(), handlers.UpdateVolunteerHandler)
		volunteers.DELETE("/:id", middleware.ValidateVolunteerMiddleware(), handlers.DeleteVolunteerHandler)
	}
}

