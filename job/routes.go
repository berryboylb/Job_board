package job

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}
var everybody []models.RoleAllowed = []models.RoleAllowed{models.PosterRole, models.AdminRole, models.SuperAdminRole}

func JobRoutes(superRoute *gin.RouterGroup) {
	jobRouter := superRoute.Group("/companies")

	jobRouter.Use(jwt.Middleware())
	jobRouter.POST("/", middleware.RolesMiddleware(everybody), create)
	jobRouter.GET("/", middleware.RolesMiddleware(everybody), get)
	jobRouter.GET("/:id", getSingle)
	jobRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), update)
	jobRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), delete)

	setupLevelRoutes(jobRouter.Group("/levels"))
	setupTypeRoutes(jobRouter.Group("/types"))
	setupApplicationRoutes(jobRouter.Group("/applications"))
}

func setupLevelRoutes(levelRouter *gin.RouterGroup) {
	levelRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), createLevel)
	levelRouter.GET("/", getLevel)
	levelRouter.GET("/:id", getSingleLevel)
	levelRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), updateLevel)
	levelRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), deleteLevel)
}

func setupTypeRoutes(sizesRouter *gin.RouterGroup) {
	sizesRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), createType)
	sizesRouter.GET("/", getType)
	sizesRouter.GET("/:id", getSingleType)
	sizesRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), updateType)
	sizesRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), deleteType)
}

func setupApplicationRoutes(sizesRouter *gin.RouterGroup) {
	sizesRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware([]models.RoleAllowed{models.UserRole}), createApplication)
	sizesRouter.GET("/", jwt.Middleware(), middleware.RolesMiddleware([]models.RoleAllowed{models.UserRole, models.AdminRole, models.SuperAdminRole}), getApplication)
	sizesRouter.GET("/:id", middleware.RolesMiddleware([]models.RoleAllowed{models.UserRole, models.AdminRole, models.SuperAdminRole}), getSingleType)
	sizesRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(everybody), updateType)
	sizesRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware( everybody), deleteType)
}
