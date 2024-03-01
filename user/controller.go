package user

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func User(ctx *gin.Context) {
	//i just want to get the user data

	ctx.JSON(http.StatusOK, gin.H{
		"data": "user data",
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


