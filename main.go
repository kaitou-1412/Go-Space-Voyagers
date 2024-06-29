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
	server := gin.Default()
	db := database.GetDB()
	routes.RegisterRoutes(server, db)
	server.Run()
}
