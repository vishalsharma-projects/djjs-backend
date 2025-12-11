package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// CreateVolunteerHandler handles volunteer creation
// @Summary Create a volunteer
// @Description Store volunteer details
// @Tags Volunteers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param volunteer body models.Volunteer true "Volunteer payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/volunteers [post]
func CreateVolunteerHandler(c *gin.Context) {
	var volunteer models.Volunteer
	if err := c.ShouldBindJSON(&volunteer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate volunteer input
	if err := validators.ValidateVolunteerInput(volunteer.VolunteerName, volunteer.BranchID, volunteer.EventID, volunteer.NumberOfDays); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateVolunteer(&volunteer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "volunteer created", "volunteer": volunteer})
}

// GetAllVolunteersHandler returns all volunteers
// @Summary Get all volunteers
// @Tags Volunteers
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Volunteer
// @Failure 500 {object} map[string]string
// @Router /api/volunteers [get]
func GetAllVolunteersHandler(c *gin.Context) {
	volunteers, err := services.GetAllVolunteers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, volunteers)
}

// GetVolunteerByEventID returns all volunteers linked to an event
// @Summary Get volunteers by event ID
// @Tags Volunteers
// @Security ApiKeyAuth
// @Produce json
// @Param event_id path int true "Event ID"
// @Success 200 {array} models.Volunteer
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/events/{event_id}/volunteers [get]
func GetVolunteerByEventID(c *gin.Context) {
	eventID := c.Param("event_id")

	evID, err := strconv.ParseUint(eventID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	vol, err := services.GetVolunteerByEventID(uint(evID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vol)
}

// UpdateVolunteerHandler updates volunteer fields
// @Summary Update a volunteer
// @Tags Volunteers
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Volunteer ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/volunteers/{id} [put]
func UpdateVolunteerHandler(c *gin.Context) {
	volunteer, exists := c.Get("volunteer")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "volunteer not found"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate update fields
	if err := validators.ValidateVolunteerUpdateFields(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := sanitizeVolunteerUpdates(payload)
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no valid fields provided"})
		return
	}

	if err := services.UpdateVolunteer(volunteer.(*models.Volunteer).ID, updates); err != nil {
		if errors.Is(err, services.ErrVolunteerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "volunteer updated"})
}

// DeleteVolunteerHandler deletes a volunteer
// @Summary Delete a volunteer
// @Tags Volunteers
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Volunteer ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/volunteers/{id} [delete]
func DeleteVolunteerHandler(c *gin.Context) {
	volunteer, exists := c.Get("volunteer")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "volunteer not found"})
		return
	}

	if err := services.DeleteVolunteer(volunteer.(*models.Volunteer).ID); err != nil {
		if errors.Is(err, services.ErrVolunteerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "volunteer deleted"})
}

func sanitizeVolunteerUpdates(payload map[string]interface{}) map[string]interface{} {
	allowed := map[string]struct{}{
		"volunteer_name": {},
		"contact":        {},
		"number_of_days": {},
		"seva_involved":  {},
		"mention_seva":   {},
		"updated_by":     {},
	}

	sanitized := make(map[string]interface{})
	for key, value := range payload {
		if _, ok := allowed[key]; ok {
			sanitized[key] = value
		}
	}

	return sanitized
}
