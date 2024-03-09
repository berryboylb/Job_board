package profile

import (
	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"
	"strconv"

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

	var req ProfileDto
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
	var (
		genderID                 uuid.UUID
		currentSalaryCurrencyID  uuid.UUID
		expectedSalaryCurrencyID uuid.UUID
		currentSalary            float64
		expectedSalary           float64
		err                      error
	)

	if id := ctx.Query("gender_id"); id != "" {
		if genderID, err = uuid.Parse(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if id := ctx.Query("current_salary_currency_id"); id != "" {
		if currentSalaryCurrencyID, err = uuid.Parse(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if id := ctx.Query("expected_salary_currency_id"); id != "" {
		if expectedSalaryCurrencyID, err = uuid.Parse(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if salary := ctx.Query("current_salary"); salary != "" {
		if currentSalary, err = strconv.ParseFloat(salary, 64); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if salary := ctx.Query("expected_salary"); salary != "" {
		if expectedSalary, err = strconv.ParseFloat(salary, 64); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	filter := ProfileDto{
		Bio:                      ctx.Query("bio"),
		Resume:                   ctx.Query("resume"),
		GenderID:                 genderID,
		CurrentSalary:            currentSalary,
		CurrentSalaryCurrencyID:  currentSalaryCurrencyID,
		ExpectedSalary:           expectedSalary,
		ExpectedSalaryCurrencyID: expectedSalaryCurrencyID,
	}

	profiles, total, page, perPage, err := getProfiles(filter, ctx.Query("page_size"), ctx.Query("page_number"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched profiles",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     profiles,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func GetSingleProfile(ctx *gin.Context) {

}

func UpdateProfile(ctx *gin.Context) {

}

func DeleteProfile(ctx *gin.Context) {

}

/* profile segment ends*/
