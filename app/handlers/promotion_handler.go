package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// CreatePromotionMaterialDetailsHandler godoc
// @Summary Create new Promotion Material Details
// @Description Create a new record in PromotionMaterialDetails table
// @Tags PromotionMaterialDetails
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body models.PromotionMaterialDetails true "Promotion Material Details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/promotion-material-details [post]
func CreatePromotionMaterialDetailsHandler(c *gin.Context) {
	var detail models.PromotionMaterialDetails
	if err := c.ShouldBindJSON(&detail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidatePromotionMaterialDetailsInput(detail.EventID, detail.PromotionMaterialID, detail.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreatePromotionMaterialDetails(&detail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Promotion Material Details created successfully",
		"data":    detail,
	})
}

// GetAllPromotionMaterialDetailsHandler godoc
// @Summary Get all Promotion Material Details
// @Description Retrieve all PromotionMaterialDetails records
// @Tags PromotionMaterialDetails
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/promotion-material-details [get]
func GetAllPromotionMaterialDetailsHandler(c *gin.Context) {
	details, err := services.GetAllPromotionMaterialDetails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Promotion Material Details fetched successfully",
		"data":    details,
	})
}

// GetPromotionMaterialDetailsByEventIDHandler godoc
// @Summary Get Promotion Material Details by Event ID
// @Description Get all Promotion Material Details records for a specific Event ID
// @Tags PromotionMaterialDetails
// @Security ApiKeyAuth
// @Produce json
// @Param event_id path int true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/promotion-material-details/event/{event_id} [get]
func GetPromotionMaterialDetailsByEventIDHandler(c *gin.Context) {
	eventIDParam := c.Param("event_id")
	eventID, err := strconv.ParseUint(eventIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	details, err := services.GetPromotionMaterialDetailsByEventID(uint(eventID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Promotion Material Details fetched successfully",
		"data":    details,
	})
}

// UpdatePromotionMaterialDetailsHandler godoc
// @Summary Update Promotion Material Details
// @Description Update a PromotionMaterialDetails record by ID
// @Tags PromotionMaterialDetails
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Promotion Material Details ID"
// @Param data body models.PromotionMaterialDetails true "Updated details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/promotion-material-details/{id} [put]
func UpdatePromotionMaterialDetailsHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.ParseUint(idParam, 10, 64)

	var detail models.PromotionMaterialDetails
	if err := c.ShouldBindJSON(&detail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to map for validation
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidatePromotionMaterialDetailsUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detail.ID = uint(id)

	if err := services.UpdatePromotionMaterialDetails(&detail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Promotion Material Details updated successfully",
		"data":    detail,
	})
}

// DeletePromotionMaterialDetailsHandler godoc
// @Summary Delete Promotion Material Details
// @Description Delete a record by ID from PromotionMaterialDetails
// @Tags PromotionMaterialDetails
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Promotion Material Details ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/promotion-material-details/{id} [delete]
func DeletePromotionMaterialDetailsHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.ParseUint(idParam, 10, 64)

	if err := services.DeletePromotionMaterialDetails(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion Material Details deleted successfully"})
}
