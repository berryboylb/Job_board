package salazrycurrency

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"

	"job_board/helpers"
	"job_board/models"
)

/* profile language segment  starts*/

func create(ctx *gin.Context) {
	var req SalaryCurrency
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	new := models.SalaryCurrency{
		Name: req.Name,
	}
	currency, err := createSalary(new)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created currency",
		StatusCode: http.StatusOK,
		Data:       currency,
	})
}

func get(ctx *gin.Context) {
	name := ctx.Query("name")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	currencies, total, page, perPage, err := getSalary(name, pageSize, pageNumber)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched currencies",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     currencies,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getSingle(ctx *gin.Context) {
	currencyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	gender, err := getSingleSalary(models.SalaryCurrency{ID: currencyID})
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched currency",
		StatusCode: http.StatusOK,
		Data:       gender,
	})
}

func update(ctx *gin.Context) {
	currencyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	var req SalaryCurrency
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	currency, err := updateSalary(currencyID, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated currency",
		StatusCode: http.StatusOK,
		Data:       currency,
	})
}

func delete(ctx *gin.Context) {
	currencyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = deleteSingle(currencyID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted currency",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

/* ProfileLanguage segment ends*/