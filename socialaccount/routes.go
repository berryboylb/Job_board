package socialaccount

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}

func SocialRoutes(superRoute *gin.RouterGroup) {
	socialRouter := superRoute.Group("/socials")

	socialRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), CreateSocial)
	socialRouter.GET("/", GetSocial)
	socialRouter.GET("/:id", GetSingleSocial)
	socialRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), UpdateSocial)
	socialRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), DeleteSocial)
}