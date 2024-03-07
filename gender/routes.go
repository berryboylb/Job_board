package gender

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}

func GenderRoutes(superRoute *gin.RouterGroup) {
	genderRouter := superRoute.Group("/genders")

	genderRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), createGender)
	genderRouter.GET("/", getGenders)
	genderRouter.GET("/:id", getSingleGender)
	genderRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), updateGender)
	genderRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), deleteGender)
}
