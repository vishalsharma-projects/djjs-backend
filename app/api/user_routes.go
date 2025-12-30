package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures user CRUD routes
func SetupUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.Use(middleware.AuthRequired())
	{
		users.POST("", 
			middleware.RequirePermission(models.ResourceUser, models.ActionCreate),
			handlers.CreateUserHandler)
		users.GET("", 
			middleware.RequirePermission(models.ResourceUser, models.ActionList),
			handlers.GetAllUsersHandler)
		users.GET("/search", 
			middleware.RequirePermission(models.ResourceUser, models.ActionList),
			handlers.GetUserSearchHandler)
		users.GET("/:id", 
			middleware.RequirePermission(models.ResourceUser, models.ActionRead),
			handlers.GetUserByIDHandler)
		users.PUT("/:id", 
			middleware.RequirePermission(models.ResourceUser, models.ActionUpdate),
			handlers.UpdateUserHandler)
		users.DELETE("/:id", 
			middleware.RequirePermission(models.ResourceUser, models.ActionDelete),
			handlers.DeleteUserHandler)
		users.POST("/:id/change-password", handlers.ChangePasswordHandler)
		users.POST("/:id/reset-password", handlers.ResetPasswordHandler)
	}
}


