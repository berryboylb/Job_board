package language

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}

func LanguageRoutes(superRoute *gin.RouterGroup) {
	languageRouter := superRoute.Group("/languages")
	languageRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), CreateLanguage)
	languageRouter.GET("/", GetLanguage)
	languageRouter.GET("/:id", GetSingleLanguage)
	languageRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), UpdateLanguage)
	languageRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), DeleteLanguage)

	ProficiencyRoutes(languageRouter.Group("/proficiencies"))
}

func ProficiencyRoutes(superRoute *gin.RouterGroup) {
	proficiencyRouter := superRoute.Group("/languages")

	proficiencyRouter.POST("/", jwt.Middleware(), middleware.RolesMiddleware(admins), CreateLanguageProficiency)
	proficiencyRouter.GET("/", GetLanguageProficiency)
	proficiencyRouter.GET("/:id", GetSingleLanguageProficiency)
	proficiencyRouter.PATCH("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), UpdateLanguageProficiency)
	proficiencyRouter.DELETE("/:id", jwt.Middleware(), middleware.RolesMiddleware(admins), DeleteLanguageProficiency)
}
