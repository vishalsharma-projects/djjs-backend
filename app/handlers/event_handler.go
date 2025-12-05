package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// ----------------------------------------------------
// Create Event
// ----------------------------------------------------

// CreateEventHandler godoc
// @Summary Create a new event
// @Tags Events
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param event body models.EventDetails true "Event payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events [post]
func CreateEventHandler(c *gin.Context) {
	var event models.EventDetails

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateEventInput(event.EventTypeID, event.EventCategoryID, event.StartDate, event.EndDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateEvent(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

// ----------------------------------------------------
// Get All Events
// ----------------------------------------------------

// GetAllEventsHandler godoc
// @Summary Get all events
// @Tags Events
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.EventDetails
// @Failure 500 {object} map[string]string
// @Router /api/events [get]
func GetAllEventsHandler(c *gin.Context) {
	events, err := services.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

// ----------------------------------------------------
// Search Events
// ----------------------------------------------------

// SearchEventsHandler godoc
// @Summary Search events
// @Description Search events by keyword
// @Tags Events
// @Security ApiKeyAuth
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {array} models.EventDetails
// @Failure 500 {object} map[string]string
// @Router /api/events/search [get]
func SearchEventsHandler(c *gin.Context) {
	search := c.Query("search")

	events, err := services.SearchEvents(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// ----------------------------------------------------
// Update Event
// ----------------------------------------------------

// UpdateEventHandler godoc
// @Summary Update an event
// @Tags Events
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body map[string]interface{} true "Updated fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events/{id} [put]
func UpdateEventHandler(c *gin.Context) {
	idParam := c.Param("id")
	eventID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateEventUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateEvent(uint(eventID), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// ----------------------------------------------------
// Delete Event
// ----------------------------------------------------

// DeleteEventHandler godoc
// @Summary Delete an event
// @Tags Events
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events/{id} [delete]
func DeleteEventHandler(c *gin.Context) {
	idParam := c.Param("id")
	eventID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := services.DeleteEvent(uint(eventID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
