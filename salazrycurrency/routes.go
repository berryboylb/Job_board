package salazrycurrency

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}

func CurrencyRoutes(superRoute *gin.RouterGroup) {
	currencyRouter := superRoute.Group("/currencies")
	currencyRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), create)
	currencyRouter.GET("/", get)
	currencyRouter.GET("/:id", getSingle)
	currencyRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), update)
	currencyRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), delete)
}