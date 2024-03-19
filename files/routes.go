package files

import (
	"github.com/gin-gonic/gin"
)

func FileRoutes(superRoute *gin.RouterGroup) {
	fileRouter := superRoute.Group("/upload")
	
	fileRouter.POST("/", uploadHandler)
}
