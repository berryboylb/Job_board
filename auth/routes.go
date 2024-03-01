package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"log"
)

// New registers the routes and returns the router.
func AuthRoutes(superRoute *gin.RouterGroup) {

	authRouter := superRoute.Group("/auth")
	{
		store := cookie.NewStore([]byte("secret"))
		authRouter.Use(sessions.Sessions("auth-session", store))
		authenticator, err := New()
		if err != nil {
			log.Fatalf("Failed to initialize the authenticator: %v", err)
		}
		authRouter.POST("/login", Login(authenticator))
		authRouter.GET("/login", Login(authenticator))
		authRouter.GET("/callback", Callback(authenticator))
		authRouter.GET("/authorize", IsAuthenticated, Authorize)
		authRouter.GET("/user", jwt.Middleware(), User)
		authRouter.GET("/logout", Logout)
	}

}
