package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/gin-gonic/gin"
)

// SetupRBACRoutes sets up routes with role-based access control examples
func SetupRBACRoutes(apiGroup *gin.RouterGroup, rbacHandler *handlers.RBACHandler) {
	// RBAC management routes - require authentication and super_admin role
	rbac := apiGroup.Group("/rbac")
	rbac.Use(middleware.AuthRequired()) // or AuthMiddleware() depending on which you use
	{
		// Role management - Super Admin only
		roles := rbac.Group("/roles")
		roles.Use(middleware.RequireRole(models.RoleTypeSuperAdmin))
		{
			roles.GET("", rbacHandler.ListRoles)
			roles.POST("", rbacHandler.CreateRole)
			roles.GET("/:id", rbacHandler.GetRole)
			roles.PUT("/:id", rbacHandler.UpdateRole)
			roles.DELETE("/:id", rbacHandler.DeleteRole)
		}

		// Permission management - Super Admin only
		permissions := rbac.Group("/permissions")
		permissions.Use(middleware.RequireRole(models.RoleTypeSuperAdmin))
		{
			permissions.GET("", rbacHandler.ListPermissions)
			permissions.POST("", rbacHandler.CreatePermission)
			permissions.GET("/:id", rbacHandler.GetPermission)
			permissions.DELETE("/:id", rbacHandler.DeletePermission)
		}

		// Role-Permission assignment - Super Admin only
		rolePermissions := rbac.Group("/role-permissions")
		rolePermissions.Use(middleware.RequireRole(models.RoleTypeSuperAdmin))
		{
			rolePermissions.POST("/grant", rbacHandler.GrantPermission)
			rolePermissions.POST("/revoke", rbacHandler.RevokePermission)
			rolePermissions.GET("/role/:roleId", rbacHandler.GetRolePermissions)
		}

		// User permissions - Authenticated users can check their own permissions
		rbac.GET("/my-permissions", rbacHandler.GetMyPermissions)
		rbac.POST("/check-permission", rbacHandler.CheckPermission)
	}
}

/*
// Example: Update user routes with RBAC
// Uncomment and use this pattern to protect your existing routes
func SetupUserRoutesWithRBAC(router *gin.Engine, userHandler *handlers.UserHandler) {
	users := router.Group("/api/users")
	users.Use(middleware.AuthRequired())
	{
		// List users - requires users:list permission
		users.GET("", 
			middleware.RequirePermission(models.ResourceUser, models.ActionList),
			userHandler.GetUsers)

		// Get single user - requires users:read permission
		users.GET("/:id", 
			middleware.RequirePermission(models.ResourceUser, models.ActionRead),
			userHandler.GetUserByID)

		// Create user - requires users:create permission
		users.POST("", 
			middleware.RequirePermission(models.ResourceUser, models.ActionCreate),
			userHandler.CreateUser)

		// Update user - requires users:update permission
		users.PUT("/:id", 
			middleware.RequirePermission(models.ResourceUser, models.ActionUpdate),
			userHandler.UpdateUser)

		// Delete user - requires users:delete permission
		users.DELETE("/:id", 
			middleware.RequirePermission(models.ResourceUser, models.ActionDelete),
			userHandler.DeleteUser)
	}
}

// Example: Update event routes with RBAC
func SetupEventRoutesWithRBAC(router *gin.Engine, eventHandler *handlers.EventHandler) {
	events := router.Group("/api/events")
	events.Use(middleware.AuthRequired())
	{
		// List events - requires events:list permission
		events.GET("", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionList),
			eventHandler.GetEvents)

		// Get single event - requires events:read permission
		events.GET("/:id", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionRead),
			eventHandler.GetEventByID)

		// Create event - requires events:create permission
		events.POST("", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionCreate),
			eventHandler.CreateEvent)

		// Update event - requires events:update permission
		events.PUT("/:id", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionUpdate),
			eventHandler.UpdateEvent)

		// Delete event - requires events:delete permission
		events.DELETE("/:id", 
			middleware.RequirePermission(models.ResourceEvent, models.ActionDelete),
			eventHandler.DeleteEvent)
	}
}

// Example: Update branch routes with RBAC
func SetupBranchRoutesWithRBAC(router *gin.Engine, branchHandler *handlers.BranchHandler) {
	branches := router.Group("/api/branches")
	branches.Use(middleware.AuthRequired())
	{
		// List branches - requires branches:list permission
		branches.GET("", 
			middleware.RequirePermission(models.ResourceBranch, models.ActionList),
			branchHandler.GetBranches)

		// Get single branch - requires branches:read permission
		branches.GET("/:id", 
			middleware.RequirePermission(models.ResourceBranch, models.ActionRead),
			branchHandler.GetBranchByID)

		// Create branch - requires branches:create permission (typically admin/coordinator)
		branches.POST("", 
			middleware.RequirePermission(models.ResourceBranch, models.ActionCreate),
			branchHandler.CreateBranch)

		// Update branch - requires branches:update permission
		branches.PUT("/:id", 
			middleware.RequirePermission(models.ResourceBranch, models.ActionUpdate),
			branchHandler.UpdateBranch)

		// Delete branch - requires branches:delete permission (typically admin only)
		branches.DELETE("/:id", 
			middleware.RequirePermission(models.ResourceBranch, models.ActionDelete),
			branchHandler.DeleteBranch)
	}
}
*/


