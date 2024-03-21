package gender

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"

	"job_board/helpers"
	"job_board/models"
)

/* gender segment  starts*/

func createGender(ctx *gin.Context) {
	var req Gender
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	newGender := models.Gender{
		Name: req.Name,
	}
	gender, err := create(newGender)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created gender",
		StatusCode: http.StatusOK,
		Data:       gender,
	})
}

func getGenders(ctx *gin.Context) {
	name := ctx.Query("name")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	genders, total, page, perPage, err := get(name, pageSize, pageNumber)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched genders",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     genders,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getSingleGender(ctx *gin.Context) {
	genderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	gender, err := getSingle(genderID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched gender",
		StatusCode: http.StatusOK,
		Data:       gender,
	})
}

func updateGender(ctx *gin.Context) {
	genderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	var req Gender
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	gender, err := update(genderID, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated gender",
		StatusCode: http.StatusOK,
		Data:       gender,
	})
}

func deleteGender(ctx *gin.Context) {
	genderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = deleteSingle(genderID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted gender",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

/* gender  segment ends*/
