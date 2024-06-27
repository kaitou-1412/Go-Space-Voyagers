package models

import (
	"math"

	"gorm.io/gorm"
)

type Planet struct {
	gorm.Model
	Name        string  `binding:"required" json:"name"`
	Description string  `binding:"required" json:"description"`
	Distance    int64   `binding:"required" json:"distance"`
	Radius      float64 `binding:"required" json:"radius"`
	Mass        float64 `json:"mass"`
	Type        PlanetType `binding:"required" json:"type"`
}

type PlanetType string

const (
	GasGiant    PlanetType = "gas_giant"
	Terrestrial PlanetType = "terrestrial"
)

var PlanetFilters = map[string]string{
    "id": "int",
    "name": "string",
    "description": "string",
    "distance": "int",
    "radius": "float",
    "mass": "float",
	"type": "string",
}

// GetFuelCost calculates the fuel cost required to travel to the planet with the given crew capacity.
func (planet Planet) GetFuelCost(crewCapacity int64) float64 {
	var gravity float64
	if planet.Type == GasGiant {
		gravity = 0.5 / math.Pow(float64(planet.Radius), 2)
	} else {
		gravity = float64(planet.Mass) / math.Pow(float64(planet.Radius), 2)

	}
	return float64(planet.Distance) / math.Pow(gravity, 2) * float64(crewCapacity)
}
