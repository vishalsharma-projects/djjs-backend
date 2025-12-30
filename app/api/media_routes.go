package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupMediaRoutes configures media CRUD routes
func SetupMediaRoutes(r *gin.RouterGroup) {
	media := r.Group("/event-media")
	media.Use(middleware.AuthRequired())
	{
		media.POST("", handlers.CreateEventMediaHandler)
		media.GET("", handlers.GetAllEventMediaHandler)
		media.GET("/event/:event_id", handlers.GetEventMediaByEventIDHandler)
		media.PUT("/:id", handlers.UpdateEventMediaHandler)
		media.DELETE("/:id", handlers.DeleteEventMediaHandler)
	}
}


