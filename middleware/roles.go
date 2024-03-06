package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"job_board/models"
)

func RolesMiddleware(roles []models.RoleAllowed) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("user")
		if !exists {
			c.String(http.StatusUnauthorized, "User not found in session")
			c.Abort()
			return
		}

		user, ok := value.(models.User)
		if !ok {
			c.String(http.StatusInternalServerError, "Mismatching types")
			c.Abort()
			return
		}

		fmt.Println(user, roles)
		// Convert user role to lowercase for case-insensitive comparison
		userRole := user.RoleName

		// Check if user role matches any of the allowed roles
		for _, allowedRole := range roles {
			if allowedRole == userRole {
				c.Next()
				return
			}
		}

		c.String(http.StatusUnauthorized, "You don't have the required role")
		c.Abort()
	}
}

