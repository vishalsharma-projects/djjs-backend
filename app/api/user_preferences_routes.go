package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupUserPreferencesRoutes configures user preferences routes
func SetupUserPreferencesRoutes(r *gin.RouterGroup) {
	prefs := r.Group("/user-preferences")
	prefs.Use(middleware.AuthRequired())
	{
		prefs.POST("", handlers.SaveUserPreferenceHandler)
		prefs.GET("", handlers.GetUserPreferenceHandler)
		prefs.DELETE("", handlers.DeleteUserPreferenceHandler)
	}
}

