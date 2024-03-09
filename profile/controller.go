package profile

import (
	// "fmt"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"

	"net/http"

	"job_board/helpers"
	"job_board/models"
	// "job_board/notifications"
)

/* profile segment  starts*/

func CreateProfile(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	var req CreateProfileDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	profile, err := createProfile(user.ID, models.Profile{
		UserID:                   user.ID,
		Bio:                      req.Bio,
		Resume:                   req.Resume,
		GenderID:                 req.GenderID,
		CurrentSalary:            req.CurrentSalary,
		CurrentSalaryCurrencyID:  &req.CurrentSalaryCurrencyID,
		ExpectedSalary:           req.ExpectedSalary,
		ExpectedSalaryCurrencyID: &req.ExpectedSalaryCurrencyID,
	})

	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created profile",
		StatusCode: http.StatusOK,
		Data:       profile,
	})
}

func GetProfile(ctx *gin.Context) {

}

func GetSingleProfile(ctx *gin.Context) {

}

func UpdateProfile(ctx *gin.Context) {

}

func DeleteProfile(ctx *gin.Context) {

}

/* profile segment ends*/
