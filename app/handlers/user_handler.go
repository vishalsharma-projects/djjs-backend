package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// CreateUserRequest represents the user creation payload
type CreateUserRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	ContactNumber string `json:"contact_number,omitempty"`
	Password      string `json:"password" binding:"required"`
	RoleID        uint   `json:"role_id" binding:"required"`
	BranchID      *uint  `json:"branch_id,omitempty"`
}

// UpdateUserRequest represents the user update payload
type UpdateUserRequest struct {
	Name          string `json:"name,omitempty"`
	Email         string `json:"email,omitempty"`
	ContactNumber string `json:"contact_number,omitempty"`
	RoleID        uint   `json:"role_id,omitempty"`
	BranchID      *uint  `json:"branch_id,omitempty"`
}

// CreateUserHandler godoc
// @Summary Create a new user
// @Description Create user with mandatory password. Password must meet strength requirements.
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User payload with password"
// @Success 201 {object} models.CreateUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users [post]
func CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	// Validate user input (name, email, contact, role)
	if err := validators.ValidateUserInput(req.Name, req.Email, req.ContactNumber, req.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password strength
	if err := validators.ValidatePasswordStrength(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user model
	user := models.User{
		Name:          req.Name,
		Email:         req.Email,
		ContactNumber: req.ContactNumber,
		RoleID:        req.RoleID,
		BranchID:      req.BranchID,
	}

	// Create user with provided password
	// Include BranchID if provided
	if err := services.CreateUser(&user, req.Password); err != nil {
		// Check if it's an email already exists error
		if err.Error() == "email already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email ID already exists. Please use a different email."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	response := models.CreateUserResponse{
		Message:  "User created successfully",
		User:     user,
		Password: "", // Don't return password for security
	}
	c.JSON(http.StatusCreated, response)
}

// GetAllUsersHandler godoc
// @Summary     Get all users
// @Tags        Users
// @Security    ApiKeyAuth
// @Produce     json
// @Success     200 {array} models.User
// @Failure     500 {object} map[string]string
// @Router      /api/users [get]
func GetAllUsersHandler(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserSearchHandler godoc
// @Summary     Search users by email or contact number
// @Description Retrieve users based on provided filters (email, contact number, or both).
// @Tags        Users
// @Security    ApiKeyAuth
// @Produce     json
// @Param       email           query string false "User Email"
// @Param       contact_number  query string false "User Contact Number"
// @Success     200 {array} models.User
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /api/users/search [get]
func GetUserSearchHandler(c *gin.Context) {
	email := c.Query("email")
	contact := c.Query("contact_number")

	// Validate search input
	if err := validators.ValidateSearchInput(email, contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := services.GetUserSearch(email, contact)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByIDHandler godoc
// @Summary     Get user by ID
// @Description Retrieve a single user by their ID
// @Tags        Users
// @Security    ApiKeyAuth
// @Produce     json
// @Param       id  path int true "User ID"
// @Success     200 {object} models.User
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /api/users/{id} [get]
func GetUserByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := services.GetUserByID(uint(userID))
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserHandler godoc
// @Summary Update a user
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "Updated fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [put]
func UpdateUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert UpdateUserRequest to map[string]interface{}
	updateData := make(map[string]interface{})
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.ContactNumber != "" {
		updateData["contact_number"] = req.ContactNumber
	}
	if req.RoleID != 0 {
		updateData["role_id"] = req.RoleID
	}
	if req.BranchID != nil {
		updateData["branch_id"] = req.BranchID
	}

	// Validate update fields
	if err := validators.ValidateUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateUser(uint(userID), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUserHandler godoc
// @Summary Delete a user (soft delete)
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [delete]
func DeleteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := services.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ChangePasswordHandler godoc
// @Summary Change user password
// @Description User can change their password by providing old and new password
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param passwordData body map[string]string true "Password change data (old_password, new_password, confirm_password)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id}/change-password [post]
func ChangePasswordHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var passwordData map[string]string
	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldPassword, ok := passwordData["old_password"]
	if !ok || oldPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "old_password is required"})
		return
	}

	newPassword, ok := passwordData["new_password"]
	if !ok || newPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new_password is required"})
		return
	}

	confirmPassword, ok := passwordData["confirm_password"]
	if !ok || confirmPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "confirm_password is required"})
		return
	}

	// Validate password change request
	if err := validators.ValidatePasswordChange(oldPassword, newPassword, confirmPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ChangePassword(uint(userID), oldPassword, newPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ResetPasswordHandler godoc
// @Summary Reset user password (admin only)
// @Description Admin can reset a user's password, generating a new temporary password
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.ResetPasswordResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id}/reset-password [post]
func ResetPasswordHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	newPassword, err := services.ResetPassword(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.ResetPasswordResponse{
		Message:  "Password reset successfully",
		Password: newPassword,
	}
	c.JSON(http.StatusOK, response)
}
