package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kaitou-1412/Go-Space-Voyagers/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func seedTestDB(db *gorm.DB) error {
	// Seed the database with test data
	planets := []models.Planet{
		{
			Name: "Jupiter", 
			Description: "A far away planet",
			Distance: 20,
			Radius: 9,
			Mass: 9,
			Type: models.GasGiant,
		},
		{
			Name: "Pluto", 
			Description: "A small planet",
			Distance: 50,
			Radius: 2,
			Mass: 2,
			Type: models.Terrestrial,
		},
	}
	result := db.Create(&planets)
	return result.Error
}

func setupDBandRouter() (*gin.Engine, string, *sql.DB, error){
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
        return nil, "Failed to open in-memory database: %v", nil, err
    }

	sqlDB, _ := db.DB()

	// Migrate the schema
	err = db.AutoMigrate(&models.Planet{})
	if err != nil {
        return nil, "Failed to migrate User model: %v", sqlDB, err
    }

	if err = seedTestDB(db); err!= nil {
		return nil, "Failed to insert test data: %v", sqlDB, err
	}

	// Set up Gin router with the handler
	gin.SetMode(gin.TestMode)
    router := gin.New()
	RegisterRoutes(router, db)

	return router, "", sqlDB, nil
}

func TestGetPlanets(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	// Create a test request
    w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/planets", nil)

	// Serve the request
    router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Data  []models.Planet `json:"data"`
		Status int    `json:"status"`
		Total  int `json:"total"`
		Page int `json:"page"`
		Limit int `json:"limit"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    } 

	assert.Equal(t, "Jupiter", response.Data[0].Name)
	assert.Equal(t, "Pluto", response.Data[1].Name)
	assert.Equal(t, 2, response.Total)

}

func TestGetPlanet(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	// Create a test request
    w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/planets/1", nil)

	// Serve the request
    router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Data  models.Planet `json:"data"`
		Status int    `json:"status"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    } 

	assert.Equal(t, "Jupiter", response.Data.Name)

}

func TestCreatePlanet(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	body := gin.H{
		"name": "Neptune",
		"description": "The far away gassy planet",
		"distance": 40,
		"radius": 8,
		"mass": 8,
		"type": "gas_giant",
	}
	jsonBody, _ := json.Marshal(body)


	// Create a test request
    w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/planets", bytes.NewBuffer(jsonBody))

	// Serve the request
    router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response struct {
		Planet  models.Planet `json:"planet"`
		Message  string `json:"message"`
		Status int    `json:"status"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }
	assert.Equal(t, "Neptune", response.Planet.Name)
	assert.Equal(t, "Planet created!", response.Message)

}

func TestUpdatePlanet(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	body := gin.H{
		"name": "Neptune",
		"description": "The far away gassy planet",
		"distance": 40,
		"radius": 8,
		"mass": 8,
		"type": "gas_giant",
	}
	jsonBody, _ := json.Marshal(body)


	// Create a test request
    w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/planets/2", bytes.NewBuffer(jsonBody))

	// Serve the request
    router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Message  string `json:"message"`
		Status int    `json:"status"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }
	assert.Equal(t, "Planet updated successfully!", response.Message)

}

func TestDeletePlanet(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	// Create a test request
    w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/planets/2", nil)

	// Serve the request
    router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Message  string `json:"message"`
		Status int    `json:"status"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }
	assert.Equal(t, "Planet deleted successfully!", response.Message)

}

func TestPlanetFuelCost(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	body := gin.H{"Capacity": 10}
	jsonBody, _ := json.Marshal(body)

	// Create a test request
    w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/planets/getFuelCost/2", bytes.NewBuffer(jsonBody))

	// Serve the request
    router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Data  float64 `json:"data"`
		Status int    `json:"status"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }
	assert.Equal(t, 2000.0, response.Data)

}