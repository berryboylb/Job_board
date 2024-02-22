package auth

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// New registers the routes and returns the router.
func Route(auth *Authenticator) *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// router.Static("/public", "web/static")
	// router.LoadHTMLGlob("web/template/*")

	router.GET("/", func(ctx *gin.Context) {
		// ctx.HTML(http.StatusOK, "home.html", nil)
		ctx.JSON(http.StatusOK, "Hello World!")
	})
	router.GET("/login", Login(auth))
	router.GET("/callback", Callback(auth))
	router.GET("/user", User)
	router.GET("/logout", Logout)
	router.GET("/protect", GinJWTMiddleware(), Protect)

	return router
}
