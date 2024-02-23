package main

import (
	"log"
	"net/http"
	"encoding/gob"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"

	"job_board/auth"
)

func Routes(authourize *auth.Authenticator) *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})


	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hello World!")
	})

	auth.AuthRoutes(authourize, router)

	return router
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	authenticator, err := auth.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	

	rtr := Routes(authenticator)

	log.Print("Server listening on http://localhost:3000/")
	
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}