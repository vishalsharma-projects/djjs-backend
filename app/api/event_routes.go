package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupEventRoutes configures event CRUD routes
func SetupEventRoutes(r *gin.RouterGroup) {
	events := r.Group("/events")
	events.Use(middleware.AuthMiddleware())
	{
		events.POST("", handlers.CreateEventHandler)
		events.GET("", handlers.GetAllEventsHandler)
		events.GET("/search", handlers.SearchEventsHandler)
		events.GET("/export", handlers.ExportEventsHandler)

		// Event-specific routes (must be before /:event_id to avoid conflicts)
		events.GET("/:event_id/specialguests", handlers.GetSpecialGuestByEventID)
		events.GET("/:event_id/specialguests/export", handlers.ExportSpecialGuestsByEventIDHandler)
		events.GET("/:event_id/volunteers", handlers.GetVolunteerByEventID)
		events.GET("/:event_id/volunteers/export", handlers.ExportVolunteersByEventIDHandler)
		events.GET("/:event_id/media/export", handlers.ExportEventMediaByEventIDHandler)
		events.GET("/:event_id/donations", handlers.GetDonationsByEvent)
		events.GET("/:event_id/promotion-materials", handlers.GetPromotionMaterialDetailsByEventIDHandler)

		events.GET("/:event_id", handlers.GetEventByIdHandler)
		events.GET("/:event_id/download", handlers.DownloadEventHandler)
		events.PUT("/:event_id", handlers.UpdateEventHandler)
		events.DELETE("/:event_id", handlers.DeleteEventHandler)
		events.PATCH("/:event_id/status", handlers.UpdateEventStatusHandler)

		// Draft routes
		events.POST("/draft", handlers.SaveDraftHandler)
		events.GET("/draft/latest", handlers.GetLatestDraftByUserHandler)
		events.GET("/draft/:draftId", handlers.GetDraftHandler)
	}
}

