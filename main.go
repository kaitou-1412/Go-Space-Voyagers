package main

import (
	"github.com/Deloitte-US/Go-Final-Assignment/db"
	"github.com/Deloitte-US/Go-Final-Assignment/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initiate the database connection
	db.InitDB()

	// Create a new Gin server instance with default middleware
	server := gin.Default()

	// Register the routes for the server
	routes.RegisterRoutes(server)

	// Run the server on port 8080
	server.Run(":8080")
}
