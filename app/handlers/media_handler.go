package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// CreateEventMediaHandler creates a new EventMedia record
// @Summary Create new Event Media
// @Description Create a new record in EventMedia table
// @Tags EventMedia
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body models.EventMedia true "Event Media Details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/event-media [post]
func CreateEventMediaHandler(c *gin.Context) {
	var media models.EventMedia
	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateEventMediaInput(media.EventID, media.MediaCoverageTypeID, media.CompanyName, media.FirstName, media.LastName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateEventMedia(&media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event Media created successfully",
		"data":    media,
	})
}

// GetAllEventMediaHandler retrieves all EventMedia records
// @Summary Get all Event Media
// @Description Retrieve all EventMedia records
// @Tags EventMedia
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/event-media [get]
func GetAllEventMediaHandler(c *gin.Context) {
	medias, err := services.GetAllEventMedia()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch records"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Event Media fetched successfully",
		"data":    medias,
	})
}

// GetEventMediaByEventIDHandler godoc
// @Summary Get Event Media by Event ID
// @Description Get all Event Media records for a specific Event ID
// @Tags EventMedia
// @Security ApiKeyAuth
// @Produce json
// @Param event_id path int true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/event-media/event/{event_id} [get]
func GetEventMediaByEventIDHandler(c *gin.Context) {
	eventIDParam := c.Param("event_id")
	eventID, err := strconv.ParseUint(eventIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	mediaList, err := services.GetEventMediaByEventID(uint(eventID))
	// Return empty array if no media found (not an error)
	if err != nil {
		mediaList = []models.EventMedia{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event Media fetched successfully",
		"data":    mediaList,
	})
}

// UpdateEventMediaHandler updates an existing EventMedia record
// @Summary Update Event Media
// @Description Update an EventMedia record by ID
// @Tags EventMedia
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Event Media ID"
// @Param data body models.EventMedia true "Updated details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/event-media/{id} [put]
func UpdateEventMediaHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var media models.EventMedia
	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to map for validation
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateEventMediaUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	media.ID = uint(id)

	if err := services.UpdateEventMedia(&media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event Media updated successfully",
		"data":    media,
	})
}

// DeleteEventMediaHandler deletes an EventMedia record by ID
// @Summary Delete Event Media
// @Description Delete a record by ID from EventMedia
// @Tags EventMedia
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Event Media ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/event-media/{id} [delete]
func DeleteEventMediaHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := services.DeleteEventMedia(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event Media deleted successfully"})
}
