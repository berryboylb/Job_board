package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"encoding/gob"

	"job_board/auth"
	"job_board/models"
)

func AddRoutes(superRoute *gin.RouterGroup) {
	//regster types
	gob.Register(map[string]interface{}{})
	gob.Register(auth.GoogleResponse{})
	gob.Register(auth.EmailResponse{})
	gob.Register(auth.GithubResponse{})
	gob.Register(models.User{})
	gob.Register(oauth2.Token{})
	gob.Register(models.AdminRole)
	gob.Register(models.UserRole)
	gob.Register(models.PosterRole)

	//register routes
	auth.AuthRoutes(superRoute)
}