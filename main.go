package main

import (
	"log"
	"mail-qq/routes/api"

	"mail-qq/queues"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	r := gin.Default()

	api.LoadAPIRoutes(r)
	queues.InitRabbitMQ()

	r.Run()
}
