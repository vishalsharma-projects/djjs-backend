package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupMasterRoutes configures master data routes for dropdowns
func SetupMasterRoutes(r *gin.RouterGroup) {
	master := r.Group("")
	master.Use(middleware.AuthMiddleware())
	{
		master.GET("/event-types", handlers.GetAllEventTypesHandler)
		master.GET("/event-categories", handlers.GetAllEventCategoriesHandler)
		master.GET("/countries", handlers.GetAllCountriesHandler)
		master.GET("/states", handlers.GetAllStatesHandler)
		master.GET("/countries/:country_id/states", handlers.GetStatesByCountryHandler)
		master.GET("/cities", handlers.GetAllCitiesHandler)
		master.GET("/cities/by-state", handlers.GetCitiesByStateHandler)
		master.GET("/districts", handlers.GetDistrictsHandler)
		master.GET("/districts/all", handlers.GetAllDistrictsHandler)
		master.GET("/promotion-material-types", handlers.GetAllPromotionMaterialTypesHandler)
		master.GET("/coordinators", handlers.GetCoordinatorDropdownHandler)
		master.GET("/orators", handlers.GetOratorDropdownHandler)
		master.GET("/languages", handlers.GetAllLanguagesHandler)
		master.GET("/seva-types", handlers.GetAllSevaTypesHandler)
		master.GET("/event-sub-categories", handlers.GetAllEventSubCategoriesHandler)
		master.GET("/event-sub-categories/by-category", handlers.GetEventSubCategoriesByCategoryHandler)
		master.GET("/roles", handlers.GetAllRolesHandler)
	}
}


