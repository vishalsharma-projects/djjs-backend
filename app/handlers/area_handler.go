package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/gin-gonic/gin"
)

// CreateAreaHandler godoc
// @Summary Create a new area
// @Tags Areas
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param area body models.Area true "Area payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/areas [post]
func CreateAreaHandler(c *gin.Context) {
	var area models.Area
	if err := c.ShouldBindJSON(&area); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateAreaInput(area.BranchID, area.DistrictID.String(), area.AreaName, area.AreaCoverage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateArea(&area); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Area created successfully",
		"area":    area,
	})
}

// GetAllAreasHandler godoc
// @Summary Get all areas
// @Tags Areas
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Area
// @Failure 500 {object} map[string]string
// @Router /api/areas [get]
func GetAllAreasHandler(c *gin.Context) {
	areas, err := services.GetAllAreas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, areas)
}

// GetAreaSearchHandler godoc
// @Summary Get areas by name (or all if none provided)
// @Description Retrieve an area by its name, or all areas if no name is provided.
// @Tags Areas
// @Security ApiKeyAuth
// @Produce json
// @Param area_name query string false "Area Name"
// @Success 200 {array} models.Area
// @Failure 404 {object} map[string]string
// @Router /api/areas/search [get]
func GetAreaSearchHandler(c *gin.Context) {
	areaName := c.Query("area_name")

	areas, err := services.GetAreaSearch(areaName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, areas)
}

// UpdateAreaHandler godoc
// @Summary Update an area
// @Tags Areas
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Area ID"
// @Param area body map[string]interface{} true "Updated fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/areas/{id} [put]
func UpdateAreaHandler(c *gin.Context) {
	idParam := c.Param("id")
	areaID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateAreaUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateArea(uint(areaID), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Area updated successfully"})
}

// DeleteAreaHandler godoc
// @Summary Delete an area
// @Tags Areas
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Area ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/areas/{id} [delete]
func DeleteAreaHandler(c *gin.Context) {
	idParam := c.Param("id")
	areaID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID"})
		return
	}

	if err := services.DeleteArea(uint(areaID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Area deleted successfully"})
}
