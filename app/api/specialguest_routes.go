package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupSpecialGuestRoutes configures special guest routes
func SetupSpecialGuestRoutes(r *gin.RouterGroup) {
	specialguests := r.Group("/specialguests")
	specialguests.Use(middleware.AuthRequired())
	{
		specialguests.POST("", handlers.CreateSpecialGuestHandler)
		specialguests.GET("", handlers.GetAllSpecialGuestsHandler)
		specialguests.PUT("/:id", middleware.ValidateSpecialGuestMiddleware(), handlers.UpdateSpecialGuestHandler)
		specialguests.DELETE("/:id", middleware.ValidateSpecialGuestMiddleware(), handlers.DeleteSpecialGuestHandler)
	}
}

