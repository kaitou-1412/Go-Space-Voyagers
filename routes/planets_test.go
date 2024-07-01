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

	tests := []struct {
		endpoint string
		expectedStatus int
		expectedName1 string
		expectedTotal int
	}{
		{"/planets", http.StatusOK, "Jupiter", 2},
		{"/planets?sort=radius", http.StatusOK, "Pluto", 2},
		{`/planets?filter[type]={"eq": "gas_giant"}`, http.StatusOK, "Jupiter", 1},
		{`/planets?filter[type]={"neq": "gas_giant"}`, http.StatusOK, "Pluto", 1},
		{`/planets?filter[radius]={"gt": 8}`, http.StatusOK, "Jupiter", 1},
		{`/planets?filter[radius]={"gte": 9}`, http.StatusOK, "Jupiter", 1},
		{`/planets?filter[radius]={"lte": 9}`, http.StatusOK, "Jupiter", 2},
		{`/planets?filter[radius]={"lt": 10}`, http.StatusOK, "Jupiter", 2},
		{`/planets?filter[type]={"like": "gas"}`, http.StatusOK, "Jupiter", 1},
		{`/planets?filter[type]={"in": ["gas_giant", "terrestrial"]}`, http.StatusOK, "Jupiter", 2},
		{`/planets?filter[type]={"notin": ["gas_giant", "terrestrial"]}`, http.StatusOK, "", 0},
		{`/planets?filter[type]={"like": 1}`, http.StatusBadRequest, "Jupiter", 0},
		{`/planets?page=1&limit=1`, http.StatusOK, "Jupiter", 1},
		{`/planets?page=abc&limit=abc`, http.StatusBadRequest, "", 0},
	}

	for _, test := range tests {
		// Create a test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", test.endpoint, nil)
	
		// Serve the request
		router.ServeHTTP(w, req)
		assert.Equal(t, test.expectedStatus, w.Code)
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
		if len(response.Data) > 0 {
			assert.Equal(t, test.expectedName1, response.Data[0].Name)
		}
		assert.Equal(t, test.expectedTotal, response.Total)
	}

}

func TestGetPlanet(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	tests := []struct {
		endpoint string
		expectedStatus int
		expectedName string
		expectedErrorMessage string
	}{
		{"/planets/1", http.StatusOK, "Jupiter", ""},
		{"/planets/3", http.StatusBadRequest, "", "Could not fetch planet."},
		{"/planets/abc", http.StatusBadRequest, "", "Could not parse planet id."},
	}

	for _, test := range tests {
		// Create a test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", test.endpoint, nil)

		// Serve the request
		router.ServeHTTP(w, req)
		assert.Equal(t, test.expectedStatus, w.Code)
		var response struct {
			Data  models.Planet `json:"data"`
			Message  string `json:"message"`
			Status int    `json:"status"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		} 

		if response.Data.ID != 0 {
			assert.Equal(t, test.expectedName, response.Data.Name)
		} else {
			assert.Equal(t, test.expectedErrorMessage, response.Message)
		}
	}
	
}

func TestCreatePlanet(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	tests := []struct {
		body gin.H
		expectedStatus int
		expectedName string
		expectedMessage string
	}{
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 40,
			"radius": 8,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusCreated, "Neptune", "Planet created!"},
		{gin.H{
			"nome": "Neptune",
			"description": "The far away gassy planet",
			"distance": 40,
			"radius": 8,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusBadRequest, "", "Could not parse request data."},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 4000,
			"radius": 8,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusBadRequest, "", "Distance should be between 10 and 1000."},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 400,
			"radius": 20,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusBadRequest, "", "Radius should be between 0.1 and 10."},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 400,
			"radius": 2,
			"mass": 20,
			"type": "terrestrial",
		}, http.StatusBadRequest, "", "Mass should be between 0.1 and 10."},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 400,
			"radius": 2,
			"mass": 8,
			"type": "gas",
		}, http.StatusBadRequest, "", "Invalid planet type."},
	}

	for _, test := range tests {

		jsonBody, _ := json.Marshal(test.body)

		// Create a test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/planets", bytes.NewBuffer(jsonBody))

		// Serve the request
		router.ServeHTTP(w, req)
		assert.Equal(t, test.expectedStatus, w.Code)
		var response struct {
			Planet  models.Planet `json:"planet"`
			Message  string `json:"message"`
			Status int    `json:"status"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if test.expectedStatus == http.StatusCreated {
			assert.Equal(t, test.expectedName, response.Planet.Name)
		}
		assert.Equal(t, test.expectedMessage, response.Message)
	}

}

func TestUpdatePlanet(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	tests := []struct {
		body gin.H
		expectedStatus int
		endpoint string
		expectedMessage string
	}{
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 40,
			"radius": 8,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusOK, "/planets/2", "Planet updated successfully!"},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 40,
			"radius": 8,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusBadRequest, "/planets/3", "Could not fetch planet for given id."},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 40,
			"radius": 8,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusBadRequest, "/planets/abc", "Could not parse planet id."},
		{gin.H{
			"nome": "Neptune",
			"description": "The far away gassy planet",
			"distance": 40,
			"radius": 8,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusBadRequest, "/planets/2", "Could not parse request data."},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 4000,
			"radius": 8,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusBadRequest, "/planets/2", "Distance should be between 10 and 1000."},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 400,
			"radius": 20,
			"mass": 8,
			"type": "gas_giant",
		}, http.StatusBadRequest, "/planets/2", "Radius should be between 0.1 and 10."},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 400,
			"radius": 2,
			"mass": 20,
			"type": "terrestrial",
		}, http.StatusBadRequest, "/planets/2", "Mass should be between 0.1 and 10."},
		{gin.H{
			"name": "Neptune",
			"description": "The far away gassy planet",
			"distance": 400,
			"radius": 2,
			"mass": 8,
			"type": "gas",
		}, http.StatusBadRequest, "/planets/2", "Invalid planet type."},
	}

	for _, test := range tests {
		jsonBody, _ := json.Marshal(test.body)

		// Create a test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", test.endpoint, bytes.NewBuffer(jsonBody))

		// Serve the request
		router.ServeHTTP(w, req)
		assert.Equal(t, test.expectedStatus, w.Code)
		var response struct {
			Message  string `json:"message"`
			Status int    `json:"status"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		assert.Equal(t, test.expectedMessage, response.Message)
	}
}

func TestDeletePlanet(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	tests := []struct {
		endpoint string
		expectedStatus int
		expectedMessage string
	}{
		{"/planets/2", http.StatusOK, "Planet deleted successfully!"},
		{"/planets/3", http.StatusBadRequest, "Could not fetch planet for given id."},
		{"/planets/abc", http.StatusBadRequest, "Could not parse planet id."},
	}

	for _, test := range tests {

		// Create a test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", test.endpoint, nil)

		// Serve the request
		router.ServeHTTP(w, req)
		assert.Equal(t, test.expectedStatus, w.Code)
		var response struct {
			Message  string `json:"message"`
			Status int    `json:"status"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		assert.Equal(t, test.expectedMessage, response.Message)
	}
}

func TestPlanetFuelCost(t *testing.T) {

	router, msg, sqlDB, err := setupDBandRouter()
	if err != nil {
		t.Errorf(msg, err)
	}
	if sqlDB != nil {
		defer sqlDB.Close()
	}

	tests := []struct {
		endpoint string
		body gin.H
		expectedStatus int
		expectedMessage string
		expectedData float64
	}{
		{"/planets/getFuelCost/1", gin.H{"Capacity": 10}, http.StatusOK, "", 5.248800000000001e+06},
		{"/planets/getFuelCost/2", gin.H{"Capacity": 10}, http.StatusOK, "", 2000.0},
		{"/planets/getFuelCost/2", gin.H{"cap": 10}, http.StatusBadRequest, "Could not parse request data.", 2000.0},
		{"/planets/getFuelCost/3", gin.H{"Capacity": 10}, http.StatusBadRequest, "Could not fetch planet for given id.", 0},
		{"/planets/getFuelCost/abc", gin.H{"Capacity": 10}, http.StatusBadRequest, "Could not parse planet id.", 0},
	}

	for _, test := range tests {

		jsonBody, _ := json.Marshal(test.body)

		// Create a test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", test.endpoint, bytes.NewBuffer(jsonBody))

		// Serve the request
		router.ServeHTTP(w, req)
		assert.Equal(t, test.expectedStatus, w.Code)
		var response struct {
			Data  float64 `json:"data"`
			Message  string `json:"message"`
			Status int    `json:"status"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if test.expectedStatus == http.StatusOK {
			assert.Equal(t, test.expectedData, response.Data)
		} else {
			assert.Equal(t, test.expectedMessage, response.Message)
		}
	}

}