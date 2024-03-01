package user

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

// New registers the routes and returns the router.
func UserRoutes(superRoute *gin.RouterGroup) {

	userRouter := superRoute.Group("/users")
	{
		userRouter.Use(jwt.Middleware())
		userRouter.GET("/", middleware.RolesMiddleware([]string{string(models.AdminRole)}), GetAllUsers)
		userRouter.GET("/user", User)
		userRouter.POST("/", middleware.RolesMiddleware([]string{string(models.AdminRole)}), CreateAdmin)
		userRouter.PATCH("/", UpdateUser)
		userRouter.DELETE("/", DeleteUser)
	}

}
