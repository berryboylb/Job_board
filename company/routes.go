package comapny

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}
var everybody []models.RoleAllowed = []models.RoleAllowed{models.PosterRole , models.AdminRole, models.SuperAdminRole}


func CompanyRoutes(superRoute *gin.RouterGroup) {
	companyRouter := superRoute.Group("/companies")

	companyRouter.Use(jwt.Middleware())
	companyRouter.POST("/", middleware.RolesMiddleware(everybody), create)
	companyRouter.GET("/",  middleware.RolesMiddleware(everybody), get)
	companyRouter.GET("/:id", getSingle)
	companyRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), update)
	companyRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), delete)
}