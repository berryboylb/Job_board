package ranking

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}

func RankingRoutes(superRoute *gin.RouterGroup) {
	rankingRouter := superRoute.Group("/genders")

	rankingRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), create)
	rankingRouter.GET("/", get)
	rankingRouter.GET("/:id", getSingle)
	rankingRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), update)
	rankingRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), delete)
}
