package handlers

import (
	"fmt"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RBACHandler struct {
	rbacService *services.RBACService
	db          *gorm.DB
}

func NewRBACHandler(db *gorm.DB) *RBACHandler {
	return &RBACHandler{
		rbacService: services.GetRBACService(),
		db:          db,
	}
}

// ListRoles godoc
// @Summary List all roles
// @Description Get a list of all roles in the system
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles [get]
func (h *RBACHandler) ListRoles(c *gin.Context) {
	var roles []models.Role
	if err := h.db.Find(&roles).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch roles")
		return
	}

	utils.OK(c, "Roles retrieved successfully", roles)
}

// GetRole godoc
// @Summary Get role by ID
// @Description Get detailed information about a specific role
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/rbac/roles/{id} [get]
func (h *RBACHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid role ID")
		return
	}

	var role models.Role
	if err := h.db.First(&role, roleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Role not found")
			return
		}
		utils.InternalServerError(c, "Failed to fetch role")
		return
	}

	utils.OK(c, "Role retrieved successfully", role)
}

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role in the system
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role body models.Role true "Role object"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/roles [post]
func (h *RBACHandler) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Validate role name
	if role.Name == "" {
		utils.BadRequest(c, "Role name is required")
		return
	}

	// Prevent creating super_admin role
	if role.Name == string(models.RoleTypeSuperAdmin) {
		utils.BadRequest(c, "Cannot create super_admin role. It is a system role.")
		return
	}

	// Check if role name already exists
	var existingRole models.Role
	if err := h.db.Where("name = ?", role.Name).First(&existingRole).Error; err == nil {
		utils.BadRequest(c, fmt.Sprintf("Role with name '%s' already exists", role.Name))
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.InternalServerError(c, "Failed to check existing role")
		return
	}

	if err := h.db.Create(&role).Error; err != nil {
		// Check for unique constraint violation
		if err.Error() == "UNIQUE constraint failed: roles.name" ||
			err.Error() == "duplicate key value violates unique constraint \"roles_name_key\"" {
			utils.BadRequest(c, fmt.Sprintf("Role with name '%s' already exists", role.Name))
			return
		}
		utils.InternalServerError(c, fmt.Sprintf("Failed to create role: %v", err))
		return
	}

	utils.Created(c, "Role created successfully", role)
}

// UpdateRole godoc
// @Summary Update a role
// @Description Update an existing role
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param role body models.Role true "Role object"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/rbac/roles/{id} [put]
func (h *RBACHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid role ID")
		return
	}

	var role models.Role
	if err := h.db.First(&role, roleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Role not found")
			return
		}
		utils.InternalServerError(c, "Failed to fetch role")
		return
	}

	// Prevent modifying super_admin role
	if role.Name == string(models.RoleTypeSuperAdmin) {
		utils.BadRequest(c, "Cannot modify super_admin role. It is a system role.")
		return
	}

	var updateData models.Role
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Validate role name
	if updateData.Name == "" {
		utils.BadRequest(c, "Role name is required")
		return
	}

	// Prevent changing to super_admin name
	if updateData.Name == string(models.RoleTypeSuperAdmin) {
		utils.BadRequest(c, "Cannot change role name to super_admin. It is a system role.")
		return
	}

	// Check if new name already exists (and it's not the current role)
	if updateData.Name != role.Name {
		var existingRole models.Role
		if err := h.db.Where("name = ? AND id != ?", updateData.Name, roleID).First(&existingRole).Error; err == nil {
			utils.BadRequest(c, fmt.Sprintf("Role with name '%s' already exists", updateData.Name))
			return
		} else if err != gorm.ErrRecordNotFound {
			utils.InternalServerError(c, "Failed to check existing role")
			return
		}
	}

	// Update fields
	role.Name = updateData.Name
	role.Description = updateData.Description

	if err := h.db.Save(&role).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: roles.name" ||
			err.Error() == "duplicate key value violates unique constraint \"roles_name_key\"" {
			utils.BadRequest(c, fmt.Sprintf("Role with name '%s' already exists", updateData.Name))
			return
		}
		utils.InternalServerError(c, fmt.Sprintf("Failed to update role: %v", err))
		return
	}

	utils.OK(c, "Role updated successfully", role)
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role from the system
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/rbac/roles/{id} [delete]
func (h *RBACHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid role ID")
		return
	}

	// Check if role exists and prevent deleting super_admin
	var role models.Role
	if err := h.db.First(&role, roleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Role not found")
			return
		}
		utils.InternalServerError(c, "Failed to fetch role")
		return
	}

	// Prevent deleting super_admin role
	if role.Name == string(models.RoleTypeSuperAdmin) {
		utils.BadRequest(c, "Cannot delete super_admin role. It is a system role.")
		return
	}

	// Check if role is assigned to any users
	var userCount int64
	if err := h.db.Model(&models.User{}).Where("role_id = ?", roleID).Count(&userCount).Error; err != nil {
		utils.InternalServerError(c, "Failed to check role usage")
		return
	}

	if userCount > 0 {
		utils.BadRequest(c, fmt.Sprintf("Cannot delete role. It is assigned to %d user(s). Please reassign users first.", userCount))
		return
	}

	if err := h.db.Delete(&role).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("Failed to delete role: %v", err))
		return
	}

	// Invalidate cache for this role
	h.rbacService.RefreshCache()

	utils.OK(c, "Role deleted successfully", nil)
}

// ListPermissions godoc
// @Summary List all permissions
// @Description Get a list of all permissions in the system
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions [get]
func (h *RBACHandler) ListPermissions(c *gin.Context) {
	var permissions []models.Permission
	if err := h.db.Find(&permissions).Error; err != nil {
		utils.InternalServerError(c, "Failed to fetch permissions")
		return
	}

	utils.OK(c, "Permissions retrieved successfully", permissions)
}

// GetPermission godoc
// @Summary Get permission by ID
// @Description Get detailed information about a specific permission
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/rbac/permissions/{id} [get]
func (h *RBACHandler) GetPermission(c *gin.Context) {
	id := c.Param("id")
	permID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid permission ID")
		return
	}

	var permission models.Permission
	if err := h.db.First(&permission, permID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Permission not found")
			return
		}
		utils.InternalServerError(c, "Failed to fetch permission")
		return
	}

	utils.OK(c, "Permission retrieved successfully", permission)
}

// CreatePermission godoc
// @Summary Create a new permission
// @Description Create a new permission in the system
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param permission body models.Permission true "Permission object"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions [post]
func (h *RBACHandler) CreatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Validate required fields
	if permission.Resource == "" {
		utils.BadRequest(c, "Resource is required")
		return
	}

	if permission.Action == "" {
		utils.BadRequest(c, "Action is required")
		return
	}

	// Validate resource type
	resourceType := models.ResourceType(permission.Resource)
	if !resourceType.IsValid() {
		utils.BadRequest(c, fmt.Sprintf("Invalid resource type: %s", permission.Resource))
		return
	}

	// Validate action type
	actionType := models.ActionType(permission.Action)
	if !actionType.IsValid() {
		utils.BadRequest(c, fmt.Sprintf("Invalid action type: %s", permission.Action))
		return
	}

	// Auto-generate name from resource and action
	permission.Name = models.PermissionString(resourceType, actionType)

	// Check if permission already exists
	var existingPermission models.Permission
	if err := h.db.Where("name = ?", permission.Name).First(&existingPermission).Error; err == nil {
		utils.BadRequest(c, fmt.Sprintf("Permission '%s' already exists", permission.Name))
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.InternalServerError(c, "Failed to check existing permission")
		return
	}

	if err := h.db.Create(&permission).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: permissions.name" ||
			err.Error() == "duplicate key value violates unique constraint \"permissions_name_key\"" {
			utils.BadRequest(c, fmt.Sprintf("Permission '%s' already exists", permission.Name))
			return
		}
		utils.InternalServerError(c, fmt.Sprintf("Failed to create permission: %v", err))
		return
	}

	utils.Created(c, "Permission created successfully", permission)
}

// DeletePermission godoc
// @Summary Delete a permission
// @Description Delete a permission from the system
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/permissions/{id} [delete]
func (h *RBACHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")
	permID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid permission ID")
		return
	}

	// Check if permission exists
	var permission models.Permission
	if err := h.db.First(&permission, permID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Permission not found")
			return
		}
		utils.InternalServerError(c, "Failed to fetch permission")
		return
	}

	// Check if permission is assigned to any roles
	var rolePermissionCount int64
	if err := h.db.Model(&models.RolePermission{}).Where("permission_id = ?", permID).Count(&rolePermissionCount).Error; err != nil {
		utils.InternalServerError(c, "Failed to check permission usage")
		return
	}

	if rolePermissionCount > 0 {
		utils.BadRequest(c, fmt.Sprintf("Cannot delete permission. It is assigned to %d role(s). Please revoke it from all roles first.", rolePermissionCount))
		return
	}

	if err := h.db.Delete(&permission).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("Failed to delete permission: %v", err))
		return
	}

	utils.OK(c, "Permission deleted successfully", nil)
}

// GrantPermissionRequest represents the request body for granting a permission
type GrantPermissionRequest struct {
	RoleID       uint `json:"role_id" binding:"required"`
	PermissionID uint `json:"permission_id" binding:"required"`
}

// GrantPermission godoc
// @Summary Grant permission to role
// @Description Grant a permission to a specific role
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GrantPermissionRequest true "Grant permission request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/role-permissions/grant [post]
func (h *RBACHandler) GrantPermission(c *gin.Context) {
	var req GrantPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Get user ID from context
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		utils.Unauthorized(c, "Unauthorized")
		return
	}

	// Get user info for audit
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		utils.InternalServerError(c, "Failed to get user")
		return
	}

	// Prevent granting permissions to super_admin (they have all permissions)
	var role models.Role
	if err := h.db.First(&role, req.RoleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, "Role not found")
			return
		}
		utils.InternalServerError(c, "Failed to get role")
		return
	}

	if role.Name == string(models.RoleTypeSuperAdmin) {
		utils.BadRequest(c, "Cannot modify permissions for super_admin role. Super admin has all permissions by default.")
		return
	}

	if err := h.rbacService.GrantPermission(req.RoleID, req.PermissionID, user.Email); err != nil {
		if err == services.ErrRoleNotFound {
			utils.BadRequest(c, "Role not found")
			return
		}
		if err == services.ErrPermissionNotFound {
			utils.BadRequest(c, "Permission not found")
			return
		}
		utils.InternalServerError(c, fmt.Sprintf("Failed to grant permission: %v", err))
		return
	}

	utils.OK(c, "Permission granted successfully", nil)
}

// RevokePermissionRequest represents the request body for revoking a permission
type RevokePermissionRequest struct {
	RoleID       uint `json:"role_id" binding:"required"`
	PermissionID uint `json:"permission_id" binding:"required"`
}

// RevokePermission godoc
// @Summary Revoke permission from role
// @Description Revoke a permission from a specific role
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RevokePermissionRequest true "Revoke permission request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/role-permissions/revoke [post]
func (h *RBACHandler) RevokePermission(c *gin.Context) {
	var req RevokePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Prevent revoking permissions from super_admin
	var role models.Role
	if err := h.db.First(&role, req.RoleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, "Role not found")
			return
		}
		utils.InternalServerError(c, "Failed to get role")
		return
	}

	if role.Name == string(models.RoleTypeSuperAdmin) {
		utils.BadRequest(c, "Cannot modify permissions for super_admin role. Super admin has all permissions by default.")
		return
	}

	if err := h.rbacService.RevokePermission(req.RoleID, req.PermissionID); err != nil {
		if err == services.ErrPermissionNotFound {
			utils.BadRequest(c, "Permission not found or not assigned to this role")
			return
		}
		utils.InternalServerError(c, fmt.Sprintf("Failed to revoke permission: %v", err))
		return
	}

	utils.OK(c, "Permission revoked successfully", nil)
}

// GetRolePermissions godoc
// @Summary Get permissions for a role
// @Description Get all permissions assigned to a specific role
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roleId path int true "Role ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/role-permissions/role/{roleId} [get]
func (h *RBACHandler) GetRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("roleId")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid role ID")
		return
	}

	permissions, err := h.rbacService.GetRolePermissions(uint(roleID))
	if err != nil {
		utils.InternalServerError(c, "Failed to get role permissions")
		return
	}

	utils.OK(c, "Permissions retrieved successfully", permissions)
}

// GetMyPermissions godoc
// @Summary Get current user's permissions
// @Description Get all permissions for the currently authenticated user
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/rbac/my-permissions [get]
func (h *RBACHandler) GetMyPermissions(c *gin.Context) {
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		utils.Unauthorized(c, "Unauthorized")
		return
	}

	permissions, err := h.rbacService.GetUserPermissions(userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get permissions")
		return
	}

	// Also get role name
	roleName, _ := c.Get("roleName")

	utils.OK(c, "Permissions retrieved successfully", gin.H{
		"permissions": permissions,
		"role":        roleName,
	})
}

// CheckPermissionRequest represents the request to check a permission
type CheckPermissionRequest struct {
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

// CheckPermission godoc
// @Summary Check if user has permission
// @Description Check if the current user has a specific permission
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CheckPermissionRequest true "Permission check request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/rbac/check-permission [post]
func (h *RBACHandler) CheckPermission(c *gin.Context) {
	var req CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Validate resource and action
	resourceType := models.ResourceType(req.Resource)
	if !resourceType.IsValid() {
		utils.BadRequest(c, fmt.Sprintf("Invalid resource type: %s", req.Resource))
		return
	}

	actionType := models.ActionType(req.Action)
	if !actionType.IsValid() {
		utils.BadRequest(c, fmt.Sprintf("Invalid action type: %s", req.Action))
		return
	}

	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		utils.Unauthorized(c, "Unauthorized")
		return
	}

	err = h.rbacService.CheckPermission(
		userID,
		models.ResourceType(req.Resource),
		models.ActionType(req.Action),
	)

	hasPermission := err == nil

	utils.OK(c, "Permission checked successfully", gin.H{
		"has_permission": hasPermission,
		"permission":     models.PermissionString(models.ResourceType(req.Resource), models.ActionType(req.Action)),
	})
}
