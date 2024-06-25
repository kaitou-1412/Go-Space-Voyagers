package models

import (
	"math"

	"github.com/Deloitte-US/Go-Final-Assignment/db"
)

type Planet struct {
	ID          int64
	Name        string  `binding:"required"`
	Description string  `binding:"required"`
	Distance    int64   `binding:"required"`
	Radius      float64 `binding:"required"`
	Mass        float64
	Type        PlanetType `binding:"required"`
}

type PlanetType string

const (
	GasGiant    PlanetType = "gas_giant"
	Terrestrial PlanetType = "terrestrial"
)

var planets = []Planet{}

// Save saves the Planet object to the database.
// It inserts a new record into the "planets" table with the values of the Planet object.
// Returns an error if there was an issue executing the SQL statement or retrieving the last inserted ID.
func (p *Planet) Save() error {
	query := `
	INSERT INTO planets(name, description, distance, radius, mass, type) 
	VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(p.Name, p.Description, p.Distance, p.Radius, p.Mass, p.Type)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	p.ID = id
	return err
}

// GetAllPlanets retrieves all the planets from the database.
// It returns a slice of Planet objects and an error, if any.
func GetAllPlanets() ([]Planet, error) {
	query := "SELECT * FROM planets"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var planets []Planet

	for rows.Next() {
		var planet Planet
		err := rows.Scan(&planet.ID, &planet.Name, &planet.Description, &planet.Distance, &planet.Radius, &planet.Mass, &planet.Type)

		if err != nil {
			return nil, err
		}

		planets = append(planets, planet)
	}

	return planets, nil
}

// GetPlanetByID retrieves a planet from the database based on the provided ID.
// It returns a pointer to the retrieved Planet struct and an error, if any.
func GetPlanetByID(id int64) (*Planet, error) {
	query := "SELECT * FROM planets WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var planet Planet
	err := row.Scan(&planet.ID, &planet.Name, &planet.Description, &planet.Distance, &planet.Radius, &planet.Mass, &planet.Type)
	if err != nil {
		return nil, err
	}

	return &planet, nil
}

// Update updates the planet record in the database with the provided values.
// It updates the name, description, distance, radius, mass, and type of the planet.
// The update is performed based on the planet's ID.
func (planet Planet) Update() error {
	query := `
	UPDATE planets
	SET name = ?, description = ?, distance = ?, radius = ?, mass = ?, type = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(planet.Name, planet.Description, planet.Distance, planet.Radius, planet.Mass, planet.Type, planet.ID)
	return err
}

// Delete deletes the planet from the database.
func (planet Planet) Delete() error {
	query := "DELETE FROM planets WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(planet.ID)
	return err
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
