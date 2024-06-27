package database

import (
	"log"

	"github.com/kaitou-1412/Go-Space-Voyagers/models"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect to database.")
	}

	err = DB.AutoMigrate(&models.Planet{}) 

	if err != nil {
		log.Fatal("Migration failure.")
	}
}

func GetDB() (*gorm.DB) {
	return DB
}
