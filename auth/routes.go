package auth

import (
	
	"github.com/gin-gonic/gin"

	"log"
)

// New registers the routes and returns the router.
func AuthRoutes(superRoute *gin.RouterGroup) {

	authRouter := superRoute.Group("/auth")
	{
		
		authenticator, err := New()
		if err != nil {
			log.Fatalf("Failed to initialize the authenticator: %v", err)
		}
		authRouter.GET("/login", Login(authenticator))
		authRouter.GET("/callback", Callback(authenticator))
		authRouter.GET("/authorize", IsAuthenticated, Authorize)
		authRouter.GET("/logout", Logout)
	}

}
