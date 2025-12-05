package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// CreateDonation godoc
// @Summary Create a new donation
// @Tags Donations
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param donation body models.Donation true "Donation Payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/donations [post]
func CreateDonation(c *gin.Context) {
	var donation models.Donation

	if err := c.ShouldBindJSON(&donation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateDonationInput(donation.EventID, donation.BranchID, donation.DonationType, donation.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateDonation(&donation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Donation created successfully",
		"donation": donation,
	})
}

// GetAllDonations godoc
// @Summary Get all donations
// @Tags Donations
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Donation
// @Failure 500 {object} map[string]string
// @Router /api/donations [get]
func GetAllDonations(c *gin.Context) {
	donations, err := services.GetAllDonations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, donations)
}

// GetDonationsByEvent godoc
// @Summary Get Donations by Event ID
// @Description Get all donation records for a specific Event ID
// @Tags Donations
// @Security ApiKeyAuth
// @Produce json
// @Param event_id path int true "Event ID"
// @Success 200 {array} models.Donation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/events/{event_id}/donations [get]
func GetDonationsByEvent(c *gin.Context) {
	eventIDParam := c.Param("event_id")

	eventID, err := strconv.ParseUint(eventIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	donations, err := services.GetDonationsByEvent(uint(eventID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, donations)
}

// UpdateDonation godoc
// @Summary Update donation
// @Tags Donations
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Donation ID"
// @Param donation body map[string]interface{} true "Updated fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/donations/{id} [put]
func UpdateDonation(c *gin.Context) {
	idStr := c.Param("id")

	donationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateDonationUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateDonation(uint(donationID), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donation updated successfully"})
}

// DeleteDonation godoc
// @Summary Delete donation
// @Tags Donations
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Donation ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/donations/{id} [delete]
func DeleteDonation(c *gin.Context) {
	idStr := c.Param("id")

	donationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	if err := services.DeleteDonation(uint(donationID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donation deleted successfully"})
}
