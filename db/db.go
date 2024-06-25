package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to database.")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

// createTables creates the "planets" table in the database if it doesn't exist.
func createTables() {
	createPlanetsTable := `
		CREATE TABLE IF NOT EXISTS planets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			distance INTEGER CHECK (distance > 10 AND distance < 1000) NOT NULL,
			radius REAL CHECK (radius > 0.1 AND radius < 10) NOT NULL,
			mass REAL CHECK (mass > 0.1 AND mass < 10) NOT NULL,
			type TEXT CHECK (type IN ('gas_giant', 'terrestrial')) NOT NULL
		)
		`

	_, err := DB.Exec(createPlanetsTable)

	if err != nil {
		panic("Could not create planets table.")
	}
}
