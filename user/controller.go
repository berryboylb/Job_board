package user

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"job_board/models"
)

func User(ctx *gin.Context) {
	//i just want to get the user data

	value, exists := ctx.Get("user")
	if !exists {
		ctx.String(http.StatusUnauthorized, "User not found in session")
		return
	}

	user, ok := value.(models.User)
	if !ok {
		ctx.String(http.StatusInternalServerError, "Mismatching types")
		ctx.Abort()
		return
	}
	// Now 'user' contains the user if it exists, and you can proceed with further processing

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "successfully fetched user",
		"data":       user,
		"statusCode": http.StatusOK,
	})
}

func GetAllUsers(ctx *gin.Context) {
	//i just want to get the user data

	ctx.JSON(http.StatusOK, gin.H{
		"data": "user data",
	})
}

func CreateAdmin(ctx *gin.Context) {
	//i just want to get the user data

	ctx.JSON(http.StatusOK, gin.H{
		"data": "user data",
	})
}

func UpdateUser(ctx *gin.Context) {
	//i just want to get the user data

	ctx.JSON(http.StatusOK, gin.H{
		"data": "user data",
	})
}

func DeleteUser(ctx *gin.Context) {
	//i just want to get the user data

	ctx.JSON(http.StatusOK, gin.H{
		"data": "user data",
	})
}
