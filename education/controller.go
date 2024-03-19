package education

import (
	// "fmt"

	"errors"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"

	// "io"
	// "log"
	"net/http"

	"job_board/helpers"
	"job_board/models"
	// "job_board/notifications"
)

/* education segment  starts*/

func CreateEducation(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	if user.Profile == nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    errors.New("please create a profile first").Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	err = req.ValidateDatesAndIsCurrent()
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	startDate, endDate, err := req.ValidateDates()
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	
	newEducation := models.Education{
		ProfileID:  user.Profile.ID,
		InstitutionName: req.InstitutionName,
		StartDate:  *startDate,
		EndDate:    endDate,
		IsCurrent:   req.IsCurrent,
		FieldOFStudy: req.FieldOFStudy,
		AcademicRankingID: req.AcademicRankingID,
		DegreeID: req.DegreeID,
		GraduationYear: req.GraduationYear,
	}

	resp, err := createEducation(newEducation)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
	}
	
	helpers.CreateResponse(ctx, helpers.Response{
		Message: "Successfully created education",
		StatusCode: http.StatusOK,
		Data: resp,
	})
}

func GetEducation(ctx *gin.Context) {

}

func GetSingleEducation(ctx *gin.Context) {

}

func UpdateEducation(ctx *gin.Context) {

}

func DeleteEducation(ctx *gin.Context) {

}

/* Education segment ends*/
