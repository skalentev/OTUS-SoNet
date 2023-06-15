package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"os"
	"otus-sonet/internal/models"
	"otus-sonet/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env not loaded!", err)
	}

	fmt.Println("Started. Hostname:", os.Getenv("HOSTNAME"))

	var RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS"),
	})
	// Create a new gin instance
	//r := gin.Default()
	r := gin.New()
	r.ForwardedByClientIP = true
	if err := r.SetTrustedProxies([]string{"10.0.0.0/8", "100.64.0.0/10", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.1"}); err != nil {
		fmt.Println("Proxies set: ", err)
	}

	models.Prom = models.NewPrometheus("http")
	models.Prom.Use(r)

	// Init Master DB
	if err := models.DB.Init(models.GetDBConfig("DB_")); err != nil {
		panic(err)
	}
	// Init ReadOnly DB
	if err := models.DBSlave.Init(models.GetDBConfig("RODB_")); err != nil {
		fmt.Println("DB slave error, using master fo RO requests. error:", err)
		models.DBSlave = models.DB
	}

	defer func() {
		if err := models.DB.Close(); err != nil {
			fmt.Println("DB close error")
		}
		if err := models.DBSlave.Close(); err != nil {
			fmt.Println("DBSlave close error")
		}
	}()

	models.Posts.Init(RedisClient, models.DB.DB)
	models.Friends.Init(RedisClient, models.DB.DB)

	go models.QueueSubscribe(context.Background(), RedisClient, "Queue")

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
