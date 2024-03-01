package user

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
)

// New registers the routes and returns the router.
func UserRoutes(superRoute *gin.RouterGroup) {

	userRouter := superRoute.Group("/users")
	{
		userRouter.Use(jwt.Middleware())
		userRouter.GET("/", GetAllUsers)
		userRouter.GET("/user", User)
		userRouter.POST("/", CreateAdmin)
		userRouter.PATCH("/", UpdateUser)
		userRouter.DELETE("/", DeleteUser)
	}

}
