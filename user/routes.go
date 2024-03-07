package user

import (
	"github.com/gin-gonic/gin"

	"job_board/jwt"
	"job_board/middleware"
	"job_board/models"
	"job_board/award"
	"job_board/education"
	"job_board/internship"
	"job_board/project"
	"job_board/work"
	"job_board/profile"
	"job_board/language"
	"job_board/socialaccount"
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
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), profile.CreateProfile)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), profile.GetProfile)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), profile.GetSingleProfile)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), profile.GetProfile)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), profile.DeleteProfile)

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
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), education.CreateEducation)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), education.GetEducation)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), education.GetSingleEducation)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), education.UpdateEducation)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), education.DeleteEducation)
}

func SetupInternshipExperienceRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), internship.CreateInternShipExperience)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), internship.GetInternShipExperience)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), internship.GetInternShipExperience)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), internship.UpdateInternShipExperience)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), internship.DeleteInternShipExperience)
}

func SetupProjectExperienceRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), project.CreateProjectExperience)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), project.GetProjectExperience)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), project.GetProjectExperience)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), project.UpdateProjectExperience)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), project.DeleteProjectExperience)
}

func SetupWorkSampleRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), work.CreateWorkSample)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), work.GetWorkSample)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), work.GetWorkSample)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), work.UpdateWorkSample)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), work.DeleteWorkSample)
}

func SetupAwardRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), award.CreateAward)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), award.GetAward)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), award.GetSingleAward)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), award.UpdateAward)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), award.DeleteAward)
}

func SetupProfileLanguageRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), language.CreateProfileLanguage)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), language.GetProfileLanguage)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), language.GetProfileLanguage)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), language.UpdateProfileLanguage)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), language.DeleteProfileLanguage)
}

func SetupSocialMediaRoutes(profileRouter *gin.RouterGroup) {
	profileRouter.Use(jwt.Middleware())
	profileRouter.POST("/", middleware.RolesMiddleware(everybody), socialaccount.CreateSocialMedia)
	profileRouter.GET("/", middleware.RolesMiddleware(admins), socialaccount.GetSocialMedia)
	profileRouter.GET("/:id", middleware.RolesMiddleware(everybody), socialaccount.GetSocialMedia)
	profileRouter.PATCH("/:id", middleware.RolesMiddleware(everybody), socialaccount.UpdateSocialMedia)
	profileRouter.DELETE("/:id", middleware.RolesMiddleware(everybody), socialaccount.DeleteSocialMedia)
}
