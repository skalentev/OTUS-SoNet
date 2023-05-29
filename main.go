package main

import (
	"github.com/gin-gonic/gin"
	"otus-sonet/models"
	"otus-sonet/routes"
)

func main() {
	// Create a new gin instance
	//r := gin.Default()
	r := gin.New()
	r.ForwardedByClientIP = true
	if err := r.SetTrustedProxies([]string{"10.0.0.0/8", "100.64.0.0/10", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.1"}); err != nil {
		panic(err)
	}

	p := models.NewPrometheus("http")
	p.Use(r)

	// Load .env file
	dbConfig := models.InitConfig()

	// Initialize DB
	models.InitDB(dbConfig)
	defer models.CloseDB()

	// Load the routes
	routes.Route(r)

	// Run the server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}
