package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupPromotionRoutes configures promotion material routes
func SetupPromotionRoutes(r *gin.RouterGroup) {
	promotion := r.Group("/promotion-material-details")
	promotion.Use(middleware.AuthRequired())
	{
		promotion.POST("", handlers.CreatePromotionMaterialDetailsHandler)
		promotion.GET("", handlers.GetAllPromotionMaterialDetailsHandler)
		promotion.GET("/event/:event_id", handlers.GetPromotionMaterialDetailsByEventIDHandler)
		promotion.PUT("/:id", handlers.UpdatePromotionMaterialDetailsHandler)
		promotion.DELETE("/:id", handlers.DeletePromotionMaterialDetailsHandler)
	}
}


