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
		userRouter.GET("/", middleware.RolesMiddleware([]string{string(models.AdminRole), string(models.SuperAdminRole)}), GetAllUsers)
		userRouter.POST("/", middleware.RolesMiddleware([]string{string(models.AdminRole), string(models.SuperAdminRole)}), CreateAdmin)
		userRouter.GET("/user", User)
		userRouter.PATCH("/user", UpdateUser)
		userRouter.DELETE("/user", DeleteUser)
		userRouter.PATCH("/user/:id", middleware.RolesMiddleware([]string{string(models.AdminRole), string(models.SuperAdminRole)}), ReinStateAccount)
	}

}
