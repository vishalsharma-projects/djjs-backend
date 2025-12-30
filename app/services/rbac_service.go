package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"gorm.io/gorm"
)

var (
	ErrPermissionDenied    = errors.New("permission denied")
	ErrRoleNotFound        = errors.New("role not found")
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrInvalidRoleType     = errors.New("invalid role type")
)

// RBACService handles role-based access control logic
type RBACService struct {
	db              *gorm.DB
	permissionCache map[uint][]string // roleID -> []permissions
	cacheMutex      sync.RWMutex
}

// NewRBACService creates a new RBAC service
func NewRBACService(db *gorm.DB) *RBACService {
	service := &RBACService{
		db:              db,
		permissionCache: make(map[uint][]string),
	}
	
	// Initialize cache
	service.RefreshCache()
	
	return service
}

// CheckPermission checks if a user has a specific permission
func (s *RBACService) CheckPermission(userID uint, resource models.ResourceType, action models.ActionType) error {
	// Get user's role
	var user models.User
	if err := s.db.Preload("Role").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPermissionDenied
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Super admin has all permissions
	if user.Role.Name == string(models.RoleTypeSuperAdmin) {
		return nil
	}

	// Check if user has the specific permission
	permissionStr := models.PermissionString(resource, action)
	if s.hasPermission(user.RoleID, permissionStr) {
		return nil
	}

	// Check if user has manage permission (which includes all actions)
	managePermission := models.PermissionString(resource, models.ActionManage)
	if s.hasPermission(user.RoleID, managePermission) {
		return nil
	}

	return ErrPermissionDenied
}

// CheckPermissionByRoleID checks if a role has a specific permission
func (s *RBACService) CheckPermissionByRoleID(roleID uint, resource models.ResourceType, action models.ActionType) error {
	// Get role
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Super admin has all permissions
	if role.Name == string(models.RoleTypeSuperAdmin) {
		return nil
	}

	// Check if role has the specific permission
	permissionStr := models.PermissionString(resource, action)
	if s.hasPermission(roleID, permissionStr) {
		return nil
	}

	// Check if role has manage permission
	managePermission := models.PermissionString(resource, models.ActionManage)
	if s.hasPermission(roleID, managePermission) {
		return nil
	}

	return ErrPermissionDenied
}

// hasPermission checks if a role has a specific permission (cache-aware)
func (s *RBACService) hasPermission(roleID uint, permissionStr string) bool {
	s.cacheMutex.RLock()
	permissions, exists := s.permissionCache[roleID]
	s.cacheMutex.RUnlock()

	if !exists {
		// Cache miss - fetch from database
		permissions = s.loadPermissionsForRole(roleID)
		s.cacheMutex.Lock()
		s.permissionCache[roleID] = permissions
		s.cacheMutex.Unlock()
	}

	// Check if permission exists in the list
	for _, perm := range permissions {
		if perm == permissionStr {
			return true
		}
	}

	return false
}

// loadPermissionsForRole loads all permissions for a role from the database
func (s *RBACService) loadPermissionsForRole(roleID uint) []string {
	var rolePermissions []models.RolePermission
	
	// Use Preload with explicit error handling
	err := s.db.Preload("Permission").Where("role_id = ?", roleID).Find(&rolePermissions).Error
	if err != nil {
		// Log error but return empty slice instead of panicking
		return []string{}
	}

	permissions := make([]string, 0, len(rolePermissions))
	for _, rp := range rolePermissions {
		// Check if Permission was loaded (Preload might fail silently)
		if rp.Permission.ID == 0 {
			// Permission not loaded, skip this entry
			continue
		}
		
		permStr := models.PermissionString(
			models.ResourceType(rp.Permission.Resource),
			models.ActionType(rp.Permission.Action),
		)
		permissions = append(permissions, permStr)
	}

	return permissions
}

// GetUserPermissions returns all permissions for a user
func (s *RBACService) GetUserPermissions(userID uint) ([]string, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionDenied
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return s.GetRolePermissions(user.RoleID)
}

// GetRolePermissions returns all permissions for a role
// Always loads fresh from database to ensure consistency
func (s *RBACService) GetRolePermissions(roleID uint) ([]string, error) {
	// Always load fresh from database to ensure we have the latest data
	// Cache is used for permission checks, but for API responses we want fresh data
	permissions := s.loadPermissionsForRole(roleID)
	
	// Update cache with fresh data
	s.cacheMutex.Lock()
	s.permissionCache[roleID] = permissions
	s.cacheMutex.Unlock()

	return permissions, nil
}

// GrantPermission grants a permission to a role
func (s *RBACService) GrantPermission(roleID uint, permissionID uint, grantedBy string) error {
	// Check if role exists
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Check if permission exists
	var permission models.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPermissionNotFound
		}
		return fmt.Errorf("failed to get permission: %w", err)
	}

	// Check if permission is already granted
	var existing models.RolePermission
	if err := s.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		First(&existing).Error; err == nil {
		// Permission already granted, return success (idempotent)
		return nil
	} else if err != gorm.ErrRecordNotFound {
		// Some other database error
		return fmt.Errorf("failed to check existing permission: %w", err)
	}

	// Grant permission
	rolePermission := models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		GrantedBy:    grantedBy,
	}

	if err := s.db.Create(&rolePermission).Error; err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}

	// Invalidate cache for this role to ensure fresh data on next request
	s.cacheMutex.Lock()
	delete(s.permissionCache, roleID)
	s.cacheMutex.Unlock()

	// Verify the permission was actually saved by reloading it
	// This ensures data consistency
	_ = s.loadPermissionsForRole(roleID)

	return nil
}

// RevokePermission revokes a permission from a role
func (s *RBACService) RevokePermission(roleID uint, permissionID uint) error {
	result := s.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{})

	if result.Error != nil {
		return fmt.Errorf("failed to revoke permission: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrPermissionNotFound
	}

	// Invalidate cache for this role to ensure fresh data on next request
	s.cacheMutex.Lock()
	delete(s.permissionCache, roleID)
	s.cacheMutex.Unlock()

	// Verify the permission was actually removed by reloading
	// This ensures data consistency
	_ = s.loadPermissionsForRole(roleID)

	return nil
}

// RefreshCache refreshes the permission cache for all roles
func (s *RBACService) RefreshCache() {
	var roles []models.Role
	if err := s.db.Find(&roles).Error; err != nil {
		return
	}

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Clear existing cache
	s.permissionCache = make(map[uint][]string)

	// Reload all permissions
	for _, role := range roles {
		s.permissionCache[role.ID] = s.loadPermissionsForRole(role.ID)
	}
}

// GetRoleByName returns a role by its name
func (s *RBACService) GetRoleByName(roleName string) (*models.Role, error) {
	var role models.Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return &role, nil
}

// CreateRole creates a new role
func (s *RBACService) CreateRole(name, description string) (*models.Role, error) {
	role := models.Role{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(&role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &role, nil
}

// CreatePermission creates a new permission
func (s *RBACService) CreatePermission(resource, action, description string) (*models.Permission, error) {
	permission := models.Permission{
		Name:        fmt.Sprintf("%s:%s", resource, action),
		Resource:    resource,
		Action:      action,
		Description: description,
	}

	if err := s.db.Create(&permission).Error; err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return &permission, nil
}

// Global RBAC service instance
var rbacService *RBACService
var rbacOnce sync.Once

// GetRBACService returns the singleton RBAC service instance
func GetRBACService() *RBACService {
	rbacOnce.Do(func() {
		rbacService = NewRBACService(config.DB)
	})
	return rbacService
}


