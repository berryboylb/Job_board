package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"job_board/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}
	models.MigrateDb()

	app := gin.New()

	router := app.Group("/api/v1")
	AddRoutes(router)

	app.Run(":3000")

	log.Print("Server listening on http://localhost:3000/")

}
