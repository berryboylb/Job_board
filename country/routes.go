package country

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}

func CountryRoutes(superRoute *gin.RouterGroup) {
	countryRouter := superRoute.Group("/countries")

	countryRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), CreateCountry)
	countryRouter.GET("/", GetCountry)
	countryRouter.GET("/:id", GetSingleCountry)
	countryRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), UpdateCountry)
	countryRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), DeleteCountry)
}