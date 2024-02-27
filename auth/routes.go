package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// New registers the routes and returns the router.
func AuthRoutes(auth *Authenticator, router *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))
	router.GET("/login", Login(auth))
	router.GET("/callback", Callback(auth))
	router.GET("/authorize", IsAuthenticated, Authorize)
	router.GET("/user", GinJWTMiddleware(), User)
	router.GET("/logout", Logout)
	router.GET("/protect", GinJWTMiddleware(), Protect)
}
