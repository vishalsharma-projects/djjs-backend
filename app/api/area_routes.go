package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupAreaRoutes configures area CRUD routes
func SetupAreaRoutes(r *gin.RouterGroup) {
	areas := r.Group("/areas")
	areas.Use(middleware.AuthRequired())
	{
		areas.POST("", handlers.CreateAreaHandler)
		areas.GET("", handlers.GetAllAreasHandler)
		areas.GET("/:id", handlers.GetAreaSearchHandler)
		areas.PUT("/:id", handlers.UpdateAreaHandler)
		areas.DELETE("/:id", handlers.DeleteAreaHandler)
	}
}


