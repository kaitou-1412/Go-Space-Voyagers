package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers the routes for handling requests.
func RegisterRoutes(server *gin.Engine, db *gorm.DB) {
	server.GET("/planets", GetPlanetsHandler(db))
	server.GET("/planets/:id", GetPlanetHandler(db))
	server.GET("/planets/getFuelCost/:id", GetFuelCostHandler(db))
	server.POST("/planets", CreatePlanetHandler(db))
	server.PUT("/planets/:id", UpdatePlanetHandler(db))
	server.DELETE("/planets/:id", DeletePlanetHandler(db))
}
