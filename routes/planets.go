package routes

import (
	"net/http"
	"strconv"

	"github.com/Deloitte-US/Go-Final-Assignment/models"
	"github.com/gin-gonic/gin"
)

// getPlanets retrieves all the planets and returns them as a JSON response.
func getPlanets(context *gin.Context) {
	planets, err := models.GetAllPlanets()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not fetch planets. Try again later."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": planets})
}

// getPlanet retrieves a planet by its ID and returns it as JSON response.
func getPlanet(context *gin.Context) {
	planetId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Could not parse planet id."})
		return
	}

	planet, err := models.GetPlanetByID(planetId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not fetch planet."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": planet})
}

// createPlanet creates a new planet based on the JSON data provided in the request body.
// It binds the JSON data to the planet model, saves it to the database, and returns the created planet as a JSON response.
func createPlanet(context *gin.Context) {
	var planet models.Planet
	err := context.ShouldBindJSON(&planet)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Could not parse request data."})
		return
	}

	if len(planet.Name) == 0 || len(planet.Description) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Name and description are required."})
		return
	}

	if planet.Type == models.GasGiant {
		planet.Mass = 5
	}

	if !(10 < planet.Distance && planet.Distance < 1000) {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Distance should be between 10 and 1000."})
		return
	}

	if !(0.1 < planet.Radius && planet.Radius < 10) {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Radius should be between 0.1 and 10."})
		return
	}

	if !(0.1 < planet.Mass && planet.Mass < 10) {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Mass should be between 0.1 and 10."})
		return
	}

	if planet.Type != models.GasGiant && planet.Type != models.Terrestrial {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid planet type."})
		return
	}

	err = planet.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not create planet. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Planet created!", "planet": planet})
}

// updatePlanet updates the details of a planet based on the provided ID.
func updatePlanet(context *gin.Context) {
	planetId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Could not parse planet id."})
		return
	}

	_, err = models.GetPlanetByID(planetId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not fetch the planet."})
		return
	}

	var updatedPlanet models.Planet
	err = context.ShouldBindJSON(&updatedPlanet)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Could not parse request data."})
		return
	}

	if len(updatedPlanet.Name) == 0 || len(updatedPlanet.Description) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Name and description are required."})
		return
	}

	if updatedPlanet.Type == models.GasGiant {
		updatedPlanet.Mass = 5
	}

	if !(10 < updatedPlanet.Distance && updatedPlanet.Distance < 1000) {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Distance should be between 10 and 1000."})
		return
	}

	if !(0.1 < updatedPlanet.Radius && updatedPlanet.Radius < 10) {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Radius should be between 0.1 and 10."})
		return
	}

	if !(0.1 < updatedPlanet.Mass && updatedPlanet.Mass < 10) {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Mass should be between 0.1 and 10."})
		return
	}

	if updatedPlanet.Type != models.GasGiant && updatedPlanet.Type != models.Terrestrial {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid planet type."})
		return
	}

	updatedPlanet.ID = planetId
	err = updatedPlanet.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not update planet."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Planet updated successfully!"})
}

// deletePlanet deletes a planet based on the provided planet ID.
func deletePlanet(context *gin.Context) {
	planetId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Could not parse planet id."})
		return
	}

	planet, err := models.GetPlanetByID(planetId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not fetch the planet."})
		return
	}

	err = planet.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not delete the planet."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Planet deleted successfully!"})
}

type Crew struct {
	Capacity int64 `binding:"required"`
}

// Function to retrieve an overall fuel cost estimation for a trip to any particular exoplanet for given crew capacity.
func getFuelCost(context *gin.Context) {
	planetId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Could not parse planet id."})
		return
	}

	planet, err := models.GetPlanetByID(planetId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Could not fetch planet."})
		return
	}

	var crew Crew
	err = context.ShouldBindJSON(&crew)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Could not parse request data."})
		return
	}

	fuelCost := planet.GetFuelCost(crew.Capacity)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": fuelCost})
}