package middleware

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/gin-gonic/gin"
)

// RequirePermission creates middleware that checks if the user has specific permission
func RequirePermission(resource models.ResourceType, action models.ActionType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware or AuthRequired)
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// Handle both uint and int64 types
		var userID uint
		switch v := userIDValue.(type) {
		case uint:
			userID = v
		case int64:
			userID = uint(v)
		case float64:
			userID = uint(v)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
			c.Abort()
			return
		}

		// Check permission
		rbacService := services.GetRBACService()
		if err := rbacService.CheckPermission(userID, resource, action); err != nil {
			if err == services.ErrPermissionDenied {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "insufficient permissions",
					"required_permission": models.PermissionString(resource, action),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates middleware that checks if the user has a specific role
func RequireRole(allowedRoles ...models.RoleType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role name from context (already set by auth middleware)
		roleNameValue, exists := c.Get("roleName")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		roleName, ok := roleNameValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid role name type"})
			c.Abort()
			return
		}

		// Check if role is in allowed roles
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if roleName == string(allowedRole) {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "insufficient role permissions",
				"required_role": allowedRoles,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission creates middleware that checks if user has ANY of the specified permissions
func RequireAnyPermission(permissions ...struct{ Resource models.ResourceType; Action models.ActionType }) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		var userID uint
		switch v := userIDValue.(type) {
		case uint:
			userID = v
		case int64:
			userID = uint(v)
		case float64:
			userID = uint(v)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
			c.Abort()
			return
		}

		rbacService := services.GetRBACService()
		hasPermission := false

		for _, perm := range permissions {
			if err := rbacService.CheckPermission(userID, perm.Resource, perm.Action); err == nil {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			requiredPerms := make([]string, len(permissions))
			for i, perm := range permissions {
				requiredPerms[i] = models.PermissionString(perm.Resource, perm.Action)
			}
			c.JSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
				"required_permissions": requiredPerms,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllPermissions creates middleware that checks if user has ALL of the specified permissions
func RequireAllPermissions(permissions ...struct{ Resource models.ResourceType; Action models.ActionType }) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		var userID uint
		switch v := userIDValue.(type) {
		case uint:
			userID = v
		case int64:
			userID = uint(v)
		case float64:
			userID = uint(v)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
			c.Abort()
			return
		}

		rbacService := services.GetRBACService()

		for _, perm := range permissions {
			if err := rbacService.CheckPermission(userID, perm.Resource, perm.Action); err != nil {
				requiredPerms := make([]string, len(permissions))
				for i, p := range permissions {
					requiredPerms[i] = models.PermissionString(p.Resource, p.Action)
				}
				c.JSON(http.StatusForbidden, gin.H{
					"error": "insufficient permissions",
					"required_permissions": requiredPerms,
					"missing_permission": models.PermissionString(perm.Resource, perm.Action),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ExtractUserID helper function to extract user ID from context
func ExtractUserID(c *gin.Context) (uint, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return 0, services.ErrPermissionDenied
	}

	switch v := userIDValue.(type) {
	case uint:
		return v, nil
	case int64:
		return uint(v), nil
	case float64:
		return uint(v), nil
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(id), nil
	default:
		return 0, services.ErrPermissionDenied
	}
}

// ExtractRoleID helper function to extract role ID from context
func ExtractRoleID(c *gin.Context) (uint, error) {
	roleIDValue, exists := c.Get("roleID")
	if !exists {
		return 0, services.ErrPermissionDenied
	}

	switch v := roleIDValue.(type) {
	case uint:
		return v, nil
	case int64:
		return uint(v), nil
	case float64:
		return uint(v), nil
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(id), nil
	default:
		return 0, services.ErrPermissionDenied
	}
}

