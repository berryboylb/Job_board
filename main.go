package main

import (
	// "context"
	"log"
	// "os"

	// apitoolkit "github.com/apitoolkit/apitoolkit-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"job_board/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}
	models.MigrateDb()

	// ctx := context.Background()
	// apitoolkitClient, err := apitoolkit.NewClient(ctx, apitoolkit.Config{APIKey: os.Getenv("API_TOOLKIT")})
	// if err != nil {
	// 	log.Fatalf("Failed to load monitoring keys: %v", err)
	// }

	app := gin.New()
	// app.Use(apitoolkitClient.GinMiddleware)
	app.MaxMultipartMemory = 8 << 20 //file size 8mb
	router := app.Group("/api/v1")

	AddRoutes(router)

	app.Run(":3000")

	log.Print("Server listening on http://localhost:3000/")

}
