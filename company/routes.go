package company

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}
var everybody []models.RoleAllowed = []models.RoleAllowed{models.PosterRole, models.AdminRole, models.SuperAdminRole}

func CompanyRoutes(superRoute *gin.RouterGroup) {
	companyRouter := superRoute.Group("/companies")

	companyRouter.Use(jwt.Middleware())
	companyRouter.POST("/", middleware.RolesMiddleware(everybody), create)
	companyRouter.GET("/", middleware.RolesMiddleware(everybody), get)
	companyRouter.GET("/:id", getSingle)
	companyRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), update)
	companyRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), delete)

	SetupIndustryRoutes(companyRouter.Group("/industries"))
	SetupSizesRoutes(companyRouter.Group("/sizes"))
}

func SetupIndustryRoutes(industryRouter *gin.RouterGroup) {
	industryRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), createIndustryHandler)
	industryRouter.GET("/", getIndustryHandler)
	industryRouter.GET("/:id", getSingleIndustryHandler)
	industryRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), updateIndustryHandler)
	industryRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), deleteIndustryHandler)
}

func SetupSizesRoutes(sizesRouter *gin.RouterGroup) {
	sizesRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), createSizes)
	sizesRouter.GET("/", getSizes)
	sizesRouter.GET("/:id", getSingleSizes)
	sizesRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), updateSizes)
	sizesRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), deleteSizes)
}
