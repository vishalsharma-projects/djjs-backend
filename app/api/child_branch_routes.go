package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupChildBranchRoutes configures child branch CRUD routes
func SetupChildBranchRoutes(r *gin.RouterGroup) {
	childBranches := r.Group("/child-branches")
	childBranches.Use(middleware.AuthRequired())
	{
		childBranches.POST("", handlers.CreateChildBranchHandler)
		childBranches.GET("", handlers.GetAllChildBranchesHandler)
		childBranches.GET("/:id", handlers.GetChildBranchHandler)
		childBranches.GET("/parent/:parent_id", handlers.GetChildBranchesByParentHandler)
		childBranches.PUT("/:id", handlers.UpdateChildBranchHandler)
		childBranches.DELETE("/:id", handlers.DeleteChildBranchHandler)

		// Child Branch Infrastructure
		childBranches.POST("/:id/infrastructure", handlers.CreateChildBranchInfrastructureHandler)
		childBranches.GET("/:id/infrastructure", handlers.GetChildBranchInfrastructureHandler)

		// Child Branch Members
		childBranches.POST("/:id/members", handlers.CreateChildBranchMemberHandler)
		childBranches.GET("/:id/members", handlers.GetChildBranchMembersHandler)
	}
}


