package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// CreateSpecialGuestHandler creates a new special guest
// @Summary Create a special guest
// @Description Store a special guest profile
// @Tags SpecialGuests
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param specialGuest body models.SpecialGuest true "Special guest payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/specialguests [post]
func CreateSpecialGuestHandler(c *gin.Context) {
	var sg models.SpecialGuest
	if err := c.ShouldBindJSON(&sg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateSpecialGuestInput(sg.EventID, sg.Prefix, sg.FirstName, sg.LastName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateSpecialGuest(&sg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Special guest created", "specialGuest": sg})
}

// GetAllSpecialGuestsHandler returns all special guests
// @Summary Get all special guests
// @Tags SpecialGuests
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.SpecialGuest
// @Failure 500 {object} map[string]string
// @Router /api/specialguests [get]
func GetAllSpecialGuestsHandler(c *gin.Context) {
	guests, err := services.GetAllSpecialGuests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, guests)
}

// GetSpecialGuestByEventID returns the special guest linked to an event
// @Summary Get special guest by event ID
// @Tags SpecialGuests
// @Security ApiKeyAuth
// @Produce json
// @Param event_id path int true "Event ID"
// @Success 200 {object} models.SpecialGuest
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/events/{event_id}/specialguests [get]
func GetSpecialGuestByEventID(c *gin.Context) {
	eventID := c.Param("event_id")

	evID, err := strconv.ParseUint(eventID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	sg, err := services.GetSpecialGuestByEventID(uint(evID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sg)
}

// ExportSpecialGuestsByEventIDHandler exports special guests to Excel
// @Summary Export special guests by event ID to Excel
// @Tags SpecialGuests
// @Security ApiKeyAuth
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param event_id path int true "Event ID"
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events/{event_id}/specialguests/export [get]
func ExportSpecialGuestsByEventIDHandler(c *gin.Context) {
	eventID := c.Param("event_id")

	evID, err := strconv.ParseUint(eventID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	guests, err := services.GetSpecialGuestByEventID(uint(evID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Get event name for filename
	event, err := services.GetEventByID(uint(evID))
	eventName := "Event"
	if err == nil && event.Theme != "" {
		eventName = event.Theme
	}

	// Export to Excel
	excelBuffer, err := services.ExportSpecialGuestsToExcel(guests, eventName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel file: " + err.Error()})
		return
	}

	filename := fmt.Sprintf("special_guests_event_%d_%s.xlsx", evID, time.Now().Format("20060102"))

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", excelBuffer.Len()))

	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBuffer.Bytes())
}

// UpdateSpecialGuestHandler updates fields of a special guest
// @Summary Update a special guest
// @Tags SpecialGuests
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Special guest ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/specialguests/{id} [put]
func UpdateSpecialGuestHandler(c *gin.Context) {
	specialGuest, exists := c.Get("specialGuest")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "special guest not found"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateSpecialGuestUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := sanitizeSpecialGuestUpdates(updateData)
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no valid fields provided"})
		return
	}

	if err := services.UpdateSpecialGuest(specialGuest.(*models.SpecialGuest).ID, updates); err != nil {
		if errors.Is(err, services.ErrSpecialGuestNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Special guest updated"})
}

// DeleteSpecialGuestHandler deletes a special guest
// @Summary Delete a special guest
// @Tags SpecialGuests
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Special guest ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/specialguests/{id} [delete]
func DeleteSpecialGuestHandler(c *gin.Context) {
	specialGuest, exists := c.Get("specialGuest")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "special guest not found"})
		return
	}

	if err := services.DeleteSpecialGuest(specialGuest.(*models.SpecialGuest).ID); err != nil {
		if errors.Is(err, services.ErrSpecialGuestNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Special guest deleted"})
}

func sanitizeSpecialGuestUpdates(payload map[string]interface{}) map[string]interface{} {
	allowed := map[string]struct{}{
		"gender":                 {},
		"prefix":                 {},
		"first_name":             {},
		"middle_name":            {},
		"last_name":              {},
		"designation":            {},
		"organization":           {},
		"email":                  {},
		"city":                   {},
		"state":                  {},
		"personal_number":        {},
		"contact_person":         {},
		"contact_person_number":  {},
		"reference_branch_id":    {},
		"reference_volunteer_id": {},
		"reference_person_name":  {},
		"updated_by":             {},
	}

	sanitized := make(map[string]interface{})
	for key, value := range payload {
		if _, ok := allowed[key]; ok {
			sanitized[key] = value
		}
	}

	return sanitized
}
