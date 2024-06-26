package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaitou-1412/Go-Space-Voyagers/db"
	"github.com/kaitou-1412/Go-Space-Voyagers/routes"
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
