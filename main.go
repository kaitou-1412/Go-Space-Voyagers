package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaitou-1412/Go-Space-Voyagers/database"
	"github.com/kaitou-1412/Go-Space-Voyagers/initialize"
	"github.com/kaitou-1412/Go-Space-Voyagers/routes"
)

func init() {
	initialize.LoadEnv()
	database.ConnectToDB()
}

func main() {
	// Create a new Gin server instance with default middleware
	server := gin.Default()

	// Register the routes for the server
	routes.RegisterRoutes(server)
	// Run the server on port 8080
	server.Run()
}
