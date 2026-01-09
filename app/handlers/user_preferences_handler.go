package handlers

import (
	"net/http"

	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	// "github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/gin-gonic/gin"
)

// SaveUserPreferenceRequest represents the request to save user preferences
type SaveUserPreferenceRequest struct {
	PreferenceType string                 `json:"preference_type" binding:"required"`
	PreferenceData map[string]interface{} `json:"preference_data" binding:"required"`
}

// SaveUserPreference godoc
// @Summary Save user preferences
// @Description Save or update user preferences (e.g., column visibility, pinned columns)
// @Tags User Preferences
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body SaveUserPreferenceRequest true "Preference data"
// @Success 200 {object} models.UserPreferences "Preferences saved successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user-preferences [post]
func SaveUserPreferenceHandler(c *gin.Context) {
	userEmail, exists := middleware.GetUserEmail(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SaveUserPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	preference, err := services.SaveUserPreference(userEmail, req.PreferenceType, req.PreferenceData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save preference"})
		return
	}

	c.JSON(http.StatusOK, preference)
}

// GetUserPreference godoc
// @Summary Get user preferences
// @Description Get user preferences by type
// @Tags User Preferences
// @Security ApiKeyAuth
// @Produce json
// @Param preference_type query string true "Preference type (e.g., 'events_list_columns')"
// @Success 200 {object} models.UserPreferences "Preferences retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Preference not found"
// @Router /api/user-preferences [get]
func GetUserPreferenceHandler(c *gin.Context) {
	userEmail, exists := middleware.GetUserEmail(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	preferenceType := c.Query("preference_type")
	if preferenceType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "preference_type is required"})
		return
	}

	preference, err := services.GetUserPreference(userEmail, preferenceType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "preference not found"})
		return
	}

	c.JSON(http.StatusOK, preference)
}

// DeleteUserPreference godoc
// @Summary Delete user preferences
// @Description Delete user preferences by type
// @Tags User Preferences
// @Security ApiKeyAuth
// @Produce json
// @Param preference_type query string true "Preference type (e.g., 'events_list_columns')"
// @Success 200 {object} map[string]string "Preferences deleted successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user-preferences [delete]
func DeleteUserPreferenceHandler(c *gin.Context) {
	userEmail, exists := middleware.GetUserEmail(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	preferenceType := c.Query("preference_type")
	if preferenceType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "preference_type is required"})
		return
	}

	if err := services.DeleteUserPreference(userEmail, preferenceType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete preference"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preference deleted successfully"})
}

