package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/gin-gonic/gin"
)

// --------------------- Event Types ---------------------

// GetAllEventTypesHandler godoc
// @Summary Get all Event Types
// @Description Returns a list of all Event Types
// @Tags EventTypes
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.EventType
// @Failure 500 {object} map[string]string
// @Router /api/event-types [get]
func GetAllEventTypesHandler(c *gin.Context) {
	list, err := services.GetAllEventTypesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --------------------- Event Categories ---------------------

// GetAllEventCategoriesHandler godoc
// @Summary Get all Event Categories
// @Description Returns a list of all Event Categories
// @Tags EventCategories
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.EventCategory
// @Failure 500 {object} map[string]string
// @Router /api/event-categories [get]
func GetAllEventCategoriesHandler(c *gin.Context) {
	list, err := services.GetAllEventCategoriesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --------------------- Countries ---------------------

// GetAllCountriesHandler godoc
// @Summary Get all countries
// @Description Returns a list of all countries
// @Tags Location
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Country
// @Failure 500 {object} map[string]string
// @Router /api/countries [get]
func GetAllCountriesHandler(c *gin.Context) {
	countries, err := services.GetAllCountriesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, countries)
}

// GetAllStatesHandler godoc
// @Summary Get all states
// @Description Returns a list of all states
// @Tags Location
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.State
// @Failure 500 {object} map[string]string
// @Router /api/states [get]
func GetAllStatesHandler(c *gin.Context) {
	states, err := services.GetAllStatesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, states)
}

// GetStatesByCountryHandler godoc
// @Summary Get states by country ID
// @Description Returns states filtered by a specific country ID
// @Tags Location
// @Security ApiKeyAuth
// @Produce json
// @Param country_id path int true "Country ID"
// @Success 200 {array} models.State
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/countries/{country_id}/states [get]
func GetStatesByCountryHandler(c *gin.Context) {
	countryIDStr := c.Param("country_id")
	countryID, err := strconv.ParseUint(countryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID"})
		return
	}

	states, err := services.GetStatesByCountryService(uint(countryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, states)
}

// --------------------- Cities ---------------------

// GetAllCitiesHandler godoc
// @Summary Get all cities
// @Description Returns a list of all cities
// @Tags Cities
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.City
// @Failure 500 {object} map[string]string
// @Router /api/cities [get]
func GetAllCitiesHandler(c *gin.Context) {
	cities, err := services.GetAllCitiesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cities)
}

// GetCitiesByStateHandler godoc
// @Summary Get cities by state
// @Description Returns cities filtered by state_id
// @Tags Cities
// @Security ApiKeyAuth
// @Produce json
// @Param state_id query int true "State ID"
// @Success 200 {array} models.City
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cities/by-state [get]
func GetCitiesByStateHandler(c *gin.Context) {
	stateIDStr := c.Query("state_id")
	if stateIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state_id query parameter is required"})
		return
	}

	stateID, err := strconv.ParseUint(stateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state_id"})
		return
	}

	cities, err := services.GetCitiesByStateService(uint(stateID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cities)
}

// --------------------- Districts ---------------------

// GetDistrictsHandler godoc
// @Summary Get districts
// @Description Returns districts filtered by state_id and/or country_id (query params optional)
// @Tags Districts
// @Security ApiKeyAuth
// @Produce json
// @Param state_id query int false "State ID"
// @Param country_id query int false "Country ID"
// @Success 200 {array} models.District
// @Failure 500 {object} map[string]string
// @Router /api/districts [get]
func GetDistrictsHandler(c *gin.Context) {
	stateIDStr := c.Query("state_id")
	countryIDStr := c.Query("country_id")

	var stateID, countryID uint
	if stateIDStr != "" {
		if id, err := strconv.ParseUint(stateIDStr, 10, 64); err == nil {
			stateID = uint(id)
		}
	}
	if countryIDStr != "" {
		if id, err := strconv.ParseUint(countryIDStr, 10, 64); err == nil {
			countryID = uint(id)
		}
	}

	districts, err := services.GetDistrictsByStateCountry(stateID, countryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, districts)
}

// GetAllDistrictsHandler godoc
// @Summary Get all districts
// @Description Returns all districts without filter
// @Tags Districts
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.District
// @Failure 500 {object} map[string]string
// @Router /api/districts/all [get]
func GetAllDistrictsHandler(c *gin.Context) {
	districts, err := services.GetAllDistricts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, districts)
}

// --------------------- Promotion Material Types ---------------------

// GetAllPromotionMaterialTypesHandler godoc
// @Summary Get all Promotion Material Types
// @Description Returns a list of all Promotion Material Types
// @Tags PromotionMaterialTypes
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.PromotionMaterialType
// @Failure 500 {object} map[string]string
// @Router /api/promotion-material-types [get]
func GetAllPromotionMaterialTypesHandler(c *gin.Context) {
	list, err := services.GetAllPromotionMaterialTypesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetCoordinatorDropdownHandler godoc
// @Summary Get Coordinator Dropdown
// @Description Returns a list of coordinators (id & name) from branch_member table where branch_role = 'coordinator'
// @Tags BranchCoordinator
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.BranchMember
// @Failure 500 {object} map[string]string
// @Router /api/coordinators [get]
func GetCoordinatorDropdownHandler(c *gin.Context) {

	list, err := services.GetCoordinatorDropdownService()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, list)
}

// GetOratorDropdownHandler godoc
// @Summary Get Orator Dropdown
// @Description Returns a list of orators (Coordinators & Preachers) with id and name from branch_member table
// @Tags Orator
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.BranchMember
// @Failure 500 {object} map[string]string
// @Router /api/orators [get]
func GetOratorDropdownHandler(c *gin.Context) {

	list, err := services.GetOratorDropdownService()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, list)
}

// --------------------- Languages ---------------------

// GetAllLanguagesHandler godoc
// @Summary Get all languages
// @Description Returns a list of all languages
// @Tags Languages
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Language
// @Failure 500 {object} map[string]string
// @Router /api/languages [get]
func GetAllLanguagesHandler(c *gin.Context) {
	languages, err := services.GetAllLanguagesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, languages)
}

// --------------------- Seva Types ---------------------

// GetAllSevaTypesHandler godoc
// @Summary Get all seva types
// @Description Returns a list of all seva types
// @Tags SevaTypes
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.SevaType
// @Failure 500 {object} map[string]string
// @Router /api/seva-types [get]
func GetAllSevaTypesHandler(c *gin.Context) {
	sevaTypes, err := services.GetAllSevaTypesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sevaTypes)
}

// --------------------- Event Sub Categories ---------------------

// GetAllEventSubCategoriesHandler godoc
// @Summary Get all event sub categories
// @Description Returns a list of all event sub categories
// @Tags EventSubCategories
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.EventSubCategory
// @Failure 500 {object} map[string]string
// @Router /api/event-sub-categories [get]
func GetAllEventSubCategoriesHandler(c *gin.Context) {
	subCategories, err := services.GetAllEventSubCategoriesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subCategories)
}

// GetEventSubCategoriesByCategoryHandler godoc
// @Summary Get event sub categories by category ID
// @Description Returns event sub categories filtered by a specific category ID
// @Tags EventSubCategories
// @Security ApiKeyAuth
// @Produce json
// @Param category_id query int true "Event Category ID"
// @Success 200 {array} models.EventSubCategory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/event-sub-categories/by-category [get]
func GetEventSubCategoriesByCategoryHandler(c *gin.Context) {
	categoryIDStr := c.Query("category_id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id query parameter is required"})
		return
	}

	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}

	subCategories, err := services.GetEventSubCategoriesByCategoryService(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subCategories)
}

// --------------------- Roles ---------------------

// GetAllRolesHandler godoc
// @Summary Get all roles
// @Description Returns a list of all roles
// @Tags Roles
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Role
// @Failure 500 {object} map[string]string
// @Router /api/roles [get]
func GetAllRolesHandler(c *gin.Context) {
	roles, err := services.GetAllRolesService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}