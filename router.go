package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"encoding/gob"

	"job_board/auth"
	"job_board/country"
	"job_board/degree"
	"job_board/files"
	"job_board/gender"
	"job_board/job"
	"job_board/language"
	"job_board/models"
	"job_board/ranking"
	"job_board/user"
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

	//register session to be used any where
	store := cookie.NewStore([]byte("secret"))
	superRoute.Use(sessions.Sessions("auth-session", store))

	//register routes
	auth.AuthRoutes(superRoute)
	user.UserRoutes(superRoute)
	gender.GenderRoutes(superRoute)
	degree.DegreeRoutes(superRoute)
	ranking.RankingRoutes(superRoute)
	language.LanguageRoutes(superRoute)
	job.JobRoutes(superRoute)
	files.FileRoutes(superRoute)
	country.CountryRoutes(superRoute)
}
