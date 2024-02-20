package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"job_board/auth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	authenticator, err := auth.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := auth.Route(authenticator)

	log.Print("Server listening on http://localhost:3000/")
	
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}