package main

import (
	"github.com/gin-gonic/gin"
	"go-auth/models"
	"go-auth/routes"
)

func main() {
	// Create a new gin instance
	r := gin.Default()

	// Load .env file
	dbConfig := models.InitConfig()

	// Initialize DB
	models.InitDB(dbConfig)
	defer models.CloseDB()

	// Load the routes
	routes.Route(r)

	// Run the server
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}

}
