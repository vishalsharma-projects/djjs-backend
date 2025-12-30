package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupBranchMediaRoutes configures branch media CRUD routes
func SetupBranchMediaRoutes(r *gin.RouterGroup) {
	media := r.Group("/branch-media")
	media.Use(middleware.AuthRequired())
	{
		media.GET("", handlers.GetAllBranchMediaHandler)
		media.GET("/branch/:branch_id", handlers.GetBranchMediaByBranchIDHandler)
	}
}

// SetupChildBranchMediaRoutes configures child branch media CRUD routes
func SetupChildBranchMediaRoutes(r *gin.RouterGroup) {
	media := r.Group("/child-branch-media")
	media.Use(middleware.AuthRequired())
	{
		media.GET("", handlers.GetAllBranchMediaHandler)
		media.GET("/branch/:branch_id", handlers.GetBranchMediaByBranchIDHandler)
	}
}


