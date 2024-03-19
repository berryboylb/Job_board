package education

import (
	// "fmt"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// "io"
	// "log"
	"net/http"
	"strconv"
	"time"

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
		ProfileID:         user.Profile.ID,
		InstitutionName:   req.InstitutionName,
		StartDate:         *startDate,
		EndDate:           endDate,
		IsCurrent:         req.IsCurrent,
		FieldOFStudy:      req.FieldOFStudy,
		AcademicRankingID: req.AcademicRankingID,
		DegreeID:          req.DegreeID,
		GraduationYear:    req.GraduationYear,
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
		Message:    "Successfully created education",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func GetEducation(ctx *gin.Context) {
	var (
		degreeID          uuid.UUID
		academicRankingID uuid.UUID
		isCurrent         bool
		graduationYear    int
		startDate         time.Time
		endDate           time.Time
		err               error
	)

	if id := ctx.Query("degree_id"); id != "" {
		if degreeID, err = uuid.Parse(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if id := ctx.Query("academic_ranking_id"); id != "" {
		if academicRankingID, err = uuid.Parse(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if boolean := ctx.Query("is_current"); boolean != "" {
		if isCurrent, err = strconv.ParseBool(boolean); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if year := ctx.Query("graduation_year"); year != "" {
		if graduationYear, err = strconv.Atoi(year); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if date := ctx.Query("start_date"); date != "" {
		if startDate, err = time.Parse("2006-01-02", date); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if date := ctx.Query("end_date"); date != "" {
		if endDate, err = time.Parse("2006-01-02", date); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}
	filter := SearchEduction{
		InstitutionName:   ctx.Query("institution_name"),
		FieldOFStudy:      ctx.Query("field_of_study"),
		DegreeID:          degreeID,
		AcademicRankingID: academicRankingID,
		GraduationYear:    graduationYear,
		StartDate:         startDate,
		EndDate:           endDate,
		IsCurrent:         isCurrent,
	}

	educations, total, page, perPage, err := getEducation(filter, ctx.Query("page_size"), ctx.Query("page_number"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched educations",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     educations,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func GetSingleEducation(ctx *gin.Context) {

}

func UpdateEducation(ctx *gin.Context) {

}

func DeleteEducation(ctx *gin.Context) {

}

/* Education segment ends*/
