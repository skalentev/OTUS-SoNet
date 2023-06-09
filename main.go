package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
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

	models.Prom = models.NewPrometheus("http")
	models.Prom.Use(r)

	// Init Master DB
	if err := models.DB.Init(models.GetDBConfig("DB_")); err != nil {
		panic(err)
	}
	// Init ReadOnly DB
	if err := models.DBRO.Init(models.GetDBConfig("RODB_")); err != nil {
		fmt.Println("DB slave error, using master fo RO requests. error:", err)
		models.DBRO = models.DB
	}

	defer func() {
		if err := models.DB.Close(); err != nil {
			fmt.Println("DB close error")
		}
		if err := models.DBRO.Close(); err != nil {
			fmt.Println("DBRO close error")
		}
	}()

	// Load the routes
	routes.Route(r)
	var addr = os.Getenv("LISTEN_ADDR")

	if addr == "" {
		addr = ":8080"
	}
	// Run the server
	if err := r.Run(addr); err != nil {
		panic(err)
	}

}
