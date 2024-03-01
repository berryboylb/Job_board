package main

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

	"job_board/auth"
	"job_board/models"
	
)

func Routes(authourize *auth.Authenticator) *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})
	gob.Register(auth.GoogleResponse{})
	gob.Register(auth.EmailResponse{})
	gob.Register(auth.GithubResponse{})
	gob.Register(oauth2.Token{})
	gob.Register(models.AdminRole)
	gob.Register(models.UserRole)
	gob.Register(models.PosterRole)

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
	// models.MigrateDb()

	authenticator, err := auth.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := Routes(authenticator)

	log.Print("Server listening on http://localhost:3000/")

	// resp, err := notifications.SendNotification(notifications.Trigger{
	// 	Name:         "berryboy",
	// 	Email:        "phemmynesce4life@gmail.com",
	// 	Title:        "test",
	// 	SubscriberID: "65d952ad34a25e4f83fc27bc",
	// 	EventID:      "welcome",
	// 	Logo:         "https://via.placeholder.com/200x200",
	// })

	// if err != nil {
	// 	log.Fatalf("Failed to send notification: %v", err)
	// }

	// log.Print("send notification", resp)

	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}

}
