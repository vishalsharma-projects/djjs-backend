package handlers

import (
	"net/http"
	"strconv"
	"time"

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
	
	// Convert to presigned URLs - fail fast on errors
	mediasWithPresignedURLs, err := services.ConvertEventMediaToPresignedURLs(c.Request.Context(), medias)
	if err != nil {
		// Fail fast - return HTTP 500 with structured error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to generate presigned URLs",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Event Media fetched successfully",
		"data":    mediasWithPresignedURLs,
	})
}

// GetEventMediaByEventIDHandler godoc
// @Summary Get Event Media by Event ID
// @Description Get Event Media records for a specific Event ID with optional cursor-based pagination
// @Tags EventMedia
// @Security ApiKeyAuth
// @Produce json
// @Param event_id path int true "Event ID"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Param cursor_created_at query string false "Cursor: created_at timestamp (RFC3339)"
// @Param cursor_id query int false "Cursor: media ID"
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

	// Parse pagination parameters
	limitParam := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 20
	}

	// Parse cursor if provided
	var cursor *services.PaginationCursor
	cursorCreatedAtStr := c.Query("cursor_created_at")
	cursorIDStr := c.Query("cursor_id")
	if cursorCreatedAtStr != "" && cursorIDStr != "" {
		cursorCreatedAt, err := time.Parse(time.RFC3339, cursorCreatedAtStr)
		if err == nil {
			cursorID, err := strconv.ParseUint(cursorIDStr, 10, 64)
			if err == nil {
				cursor = &services.PaginationCursor{
					CreatedAt: cursorCreatedAt,
					ID:        uint(cursorID),
				}
			}
		}
	}

	// Use cursor-based pagination
	paginatedResult, err := services.GetEventMediaByEventIDPaginated(uint(eventID), limit, cursor)
	if err != nil {
		// Fallback to non-paginated for backward compatibility
		mediaList, fallbackErr := services.GetEventMediaByEventID(uint(eventID))
		if fallbackErr != nil {
			mediaList = []models.EventMedia{}
		}
		// Convert to presigned URLs - fail fast on errors
		mediaListWithPresignedURLs, fallbackErr := services.ConvertEventMediaToPresignedURLs(c.Request.Context(), mediaList)
		if fallbackErr != nil {
			// Fail fast - return HTTP 500 with structured error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to generate presigned URLs",
				"details": fallbackErr.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Event Media fetched successfully",
			"data":    mediaListWithPresignedURLs,
		})
		return
	}

	// Convert to presigned URLs - fail fast on errors
	mediaListWithPresignedURLs, err := services.ConvertEventMediaToPresignedURLs(c.Request.Context(), paginatedResult.Data)
	if err != nil {
		// Fail fast - return HTTP 500 with structured error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to generate presigned URLs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Event Media fetched successfully",
		"data":       mediaListWithPresignedURLs,
		"next_cursor": paginatedResult.NextCursor,
		"has_more":   paginatedResult.HasMore,
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
