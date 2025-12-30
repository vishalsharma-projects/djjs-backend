package models

// RoleType defines standard role types in the system
type RoleType string

const (
	RoleTypeSuperAdmin  RoleType = "super_admin"
	RoleTypeAdmin       RoleType = "admin"
	RoleTypeCoordinator RoleType = "coordinator"
	RoleTypeStaff       RoleType = "staff"
)

// ResourceType defines the resources in the system
type ResourceType string

const (
	ResourceUser         ResourceType = "users"
	ResourceRole         ResourceType = "roles"
	ResourcePermission   ResourceType = "permissions"
	ResourceBranch       ResourceType = "branches"
	ResourceArea         ResourceType = "areas"
	ResourceEvent        ResourceType = "events"
	ResourceDonation     ResourceType = "donations"
	ResourceVolunteer    ResourceType = "volunteers"
	ResourceSpecialGuest ResourceType = "special_guests"
	ResourceMedia        ResourceType = "media"
	ResourcePromotion    ResourceType = "promotions"
	ResourceMaster       ResourceType = "master_data"
)

// ActionType defines the actions that can be performed
type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionRead   ActionType = "read"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionList   ActionType = "list"
	ActionManage ActionType = "manage" // Full control including special operations
)

// PermissionString generates a permission string from resource and action
func PermissionString(resource ResourceType, action ActionType) string {
	return string(resource) + ":" + string(action)
}

// IsValid checks if the role type is valid
func (r RoleType) IsValid() bool {
	switch r {
	case RoleTypeSuperAdmin, RoleTypeAdmin, RoleTypeCoordinator, RoleTypeStaff:
		return true
	}
	return false
}

// String returns the string representation of RoleType
func (r RoleType) String() string {
	return string(r)
}

// IsValid checks if the resource type is valid
func (r ResourceType) IsValid() bool {
	switch r {
	case ResourceUser, ResourceRole, ResourcePermission, ResourceBranch,
		ResourceArea, ResourceEvent, ResourceDonation, ResourceVolunteer,
		ResourceSpecialGuest, ResourceMedia, ResourcePromotion, ResourceMaster:
		return true
	}
	return false
}

// String returns the string representation of ResourceType
func (r ResourceType) String() string {
	return string(r)
}

// IsValid checks if the action type is valid
func (a ActionType) IsValid() bool {
	switch a {
	case ActionCreate, ActionRead, ActionUpdate, ActionDelete, ActionList, ActionManage:
		return true
	}
	return false
}

// String returns the string representation of ActionType
func (a ActionType) String() string {
	return string(a)
}
