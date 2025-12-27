package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupBranchRoutes configures branch CRUD routes
func SetupBranchRoutes(r *gin.RouterGroup) {
	branches := r.Group("/branches")
	branches.Use(middleware.AuthMiddleware())
	{
		branches.POST("", handlers.CreateBranchHandler)
		branches.GET("", handlers.GetAllBranchesHandler)
		branches.GET("/search", handlers.GetBranchSearchHandler)
		branches.GET("/export", handlers.ExportBranchesHandler) // Must be before /:id route
		branches.GET("/parent/:parent_id/children", handlers.GetChildBranchesHandler)
		branches.GET("/:id", handlers.GetBranchHandler)
		branches.PUT("/:id", handlers.UpdateBranchHandler)
		branches.DELETE("/:id", handlers.DeleteBranchHandler)
	}

	// Branch Infrastructure routes
	branchInfra := r.Group("/branch-infra")
	branchInfra.Use(middleware.AuthMiddleware())
	{
		branchInfra.POST("", handlers.CreateBranchInfrastructureHandler)
		branchInfra.GET("", handlers.GetAllBranchInfrastructureHandler)
		branchInfra.GET("/branch/:branch_id", handlers.GetInfrastructureByBranchHandler)
		branchInfra.PUT("/:id", handlers.UpdateBranchInfrastructureHandler)
		branchInfra.DELETE("/:id", handlers.DeleteBranchInfrastructureHandler)
	}

	// Branch Member routes
	branchMember := r.Group("/branch-member")
	branchMember.Use(middleware.AuthMiddleware())
	{
		branchMember.POST("", handlers.CreateBranchMemberHandler)
		branchMember.GET("", handlers.GetAllBranchMembersHandler)
		branchMember.GET("/export", handlers.ExportMembersHandler) // Must be before /:id route
		branchMember.GET("/branch/:branch_id", handlers.GetMembersByBranchHandler)
		branchMember.PUT("/:id", handlers.UpdateBranchMemberHandler)
		branchMember.DELETE("/:id", handlers.DeleteBranchMemberHandler)
	}
}


