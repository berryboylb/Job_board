package country

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"

	"job_board/helpers"
	"job_board/models"
)

func CreateCountry(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
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

	newCountry := models.Country{
		Name: req.Name,
	}

	resp, err := createProject(newCountry, *user)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully created project",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func GetCountry(ctx *gin.Context) {
	

	filter := Request{
		Name: ctx.Query("name"),
	}

	resp, total, page, perPage, err := getProject(filter, ctx.Query("page_size"), ctx.Query("page_number"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched countries",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     resp,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func GetSingleCountry(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	ID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	resp, err := getSingleProject(ID, *user)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully fetched country",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func UpdateCountry(ctx *gin.Context) {
	// Get user from context
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	ID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
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

	resp, err := updateProject(ID, *user, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully updated internship",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func DeleteCountry(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	ID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	err = deleteSingleProject(ID, *user)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully deleted project",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}
