package degree

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}

func DegreeRoutes(superRoute *gin.RouterGroup) {
	degreeRouter := superRoute.Group("/degrees")

	degreeRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), create)
	degreeRouter.GET("/", get)
	degreeRouter.GET("/:id", getSingle)
	degreeRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), update)
	degreeRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), delete)
}
