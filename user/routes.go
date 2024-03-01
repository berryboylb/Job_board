package user

import (
	"github.com/gin-gonic/gin"
)


// New registers the routes and returns the router.
func UserRoutes(superRoute *gin.RouterGroup) {

	userRouter := superRoute.Group("/users")
	{
		userRouter.POST("/user", User)
	}

}


