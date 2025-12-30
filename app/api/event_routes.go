package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/gin-gonic/gin"
)

// SetupEventRoutes configures event CRUD routes
func SetupEventRoutes(r *gin.RouterGroup) {
	events := r.Group("/events")
	events.Use(middleware.AuthMiddleware())
	{
		events.POST("", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionCreate),
			handlers.CreateEventHandler)
		events.GET("", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionList),
			handlers.GetAllEventsHandler)
		events.GET("/search", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionList),
			handlers.SearchEventsHandler)
		events.GET("/export", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionList),
			handlers.ExportEventsHandler)

		// Event-specific routes (must be before /:event_id to avoid conflicts)
		events.GET("/:event_id/specialguests", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.GetSpecialGuestByEventID)
		events.GET("/:event_id/specialguests/export", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.ExportSpecialGuestsByEventIDHandler)
		events.GET("/:event_id/volunteers", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.GetVolunteerByEventID)
		events.GET("/:event_id/volunteers/export", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.ExportVolunteersByEventIDHandler)
		events.GET("/:event_id/media/export", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.ExportEventMediaByEventIDHandler)
		events.GET("/:event_id/donations", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.GetDonationsByEvent)
		events.GET("/:event_id/promotion-materials", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.GetPromotionMaterialDetailsByEventIDHandler)

		events.GET("/:event_id", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.GetEventByIdHandler)
		events.GET("/:event_id/download", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.DownloadEventHandler)
		events.PUT("/:event_id", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionUpdate),
			handlers.UpdateEventHandler)
		events.DELETE("/:event_id", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionDelete),
			handlers.DeleteEventHandler)
		events.PATCH("/:event_id/status", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionUpdate),
			handlers.UpdateEventStatusHandler)

		// Draft routes
		events.POST("/draft", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionCreate),
			handlers.SaveDraftHandler)
		events.GET("/draft/latest", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.GetLatestDraftByUserHandler)
		events.GET("/draft/:draftId", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			handlers.GetDraftHandler)
	}
}

