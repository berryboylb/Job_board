package user

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
)

var admins []models.RoleAllowed = []models.RoleAllowed{models.AdminRole, models.SuperAdminRole}
var everybody []models.RoleAllowed = []models.RoleAllowed{models.UserRole, models.AdminRole, models.SuperAdminRole}

// New registers the routes and returns the router.
func UserRoutes(superRoute *gin.RouterGroup) {

	userRouter := superRoute.Group("/users")
	{
		userRouter.Use(jwt.Middleware())
		userRouter.GET("/", middleware.RolesMiddleware([]models.RoleAllowed{models.AdminRole, models.SuperAdminRole}), GetAllUsers)
		userRouter.POST("/", middleware.RolesMiddleware([]models.RoleAllowed{models.AdminRole, models.SuperAdminRole}), CreateAdmin)
		userRouter.GET("/user", User)
		userRouter.PATCH("/user", UpdateUser)
		userRouter.PATCH("/user/:id", middleware.RolesMiddleware([]models.RoleAllowed{models.AdminRole, models.SuperAdminRole}), ReinStateAccount)
		userRouter.DELETE("/user", DeleteUser)

		SetupProfileRoutes(userRouter.Group("/profiles"))
	}

}

func SetupProfileRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), CreateProfile)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), GetProfile)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), GetSingleProfile)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), GetProfile)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), DeleteProfile)

	/* subprofile routes */
	SetupEducationRoutes(profileRouter.Group("/educations"))
	SetupInternshipExperienceRoutes(profileRouter.Group("/internships"))
	SetupProjectExperienceRoutes(profileRouter.Group("/projects"))
	SetupWorkSampleRoutes(profileRouter.Group("/works"))
	SetupAwardRoutes(profileRouter.Group("/awards"))
	SetupProfileLanguageRoutes(profileRouter.Group("/languages"))
	SetupSocialMediaRoutes(profileRouter.Group("/socials"))

}

func SetupEducationRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), CreateEducation)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), GetEducation)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), GetSingleEducation)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), UpdateEducation)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), DeleteEducation)
}

func SetupInternshipExperienceRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), CreateInternShipExperience)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), GetInternShipExperience)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), GetInternShipExperience)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), UpdateInternShipExperience)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), DeleteInternShipExperience)
}

func SetupProjectExperienceRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), CreateProjectExperience)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), GetProjectExperience)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), GetProjectExperience)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), UpdateProjectExperience)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), DeleteProjectExperience)
}

func SetupWorkSampleRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), CreateWorkSample)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), GetWorkSample)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), GetWorkSample)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), UpdateWorkSample)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), DeleteWorkSample)
}

func SetupAwardRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), CreateAward)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), GetAward)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), GetAward)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), UpdateAward)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), DeleteAward)
}

func SetupProfileLanguageRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), CreateProfileLanguage)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), GetProfileLanguage)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), GetProfileLanguage)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), UpdateProfileLanguage)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), DeleteProfileLanguage)
}

func SetupSocialMediaRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), CreateSocialMedia)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), GetSocialMedia)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), GetSocialMedia)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), UpdateSocialMedia)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), DeleteSocialMedia)
}
