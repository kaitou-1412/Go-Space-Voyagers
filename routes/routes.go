package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the routes for handling requests.
func RegisterRoutes(server *gin.Engine) {
	server.GET("/planets", getPlanets)
	server.GET("/planets/:id", getPlanet)
	server.GET("/planets/getFuelCost/:id", getFuelCost)
	server.POST("/planets", createPlanet)
	server.PUT("/planets/:id", updatePlanet)
	server.DELETE("/planets/:id", deletePlanet)
}
