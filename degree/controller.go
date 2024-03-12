package degree

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"

	"job_board/helpers"
	"job_board/models"
)

/* profile language segment  starts*/

func create(ctx *gin.Context) {
	var req Degree
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	new := models.Degree{
		Name: req.Name,
	}
	degree, err := createDegree(new)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created degree",
		StatusCode: http.StatusOK,
		Data:       degree,
	})
}

func get(ctx *gin.Context) {
	name := ctx.Query("name")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	degrees, total, page, perPage, err := getDegree(name, pageSize, pageNumber)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched degrees",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     degrees,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getSingle(ctx *gin.Context) {
	degreeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	degree, err := getSingleDegree(models.Degree{ID: degreeID})
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched degree",
		StatusCode: http.StatusOK,
		Data:       degree,
	})
}

func update(ctx *gin.Context) {
	degreeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	var req Degree
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	degree, err := updateDegree(degreeID, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated degree",
		StatusCode: http.StatusOK,
		Data:       degree,
	})
}

func delete(ctx *gin.Context) {
	degreeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = deleteSingle(degreeID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted degree",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

/* ProfileLanguage segment ends*/
