package company

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"

	"job_board/helpers"
	"job_board/models"
)


func create(ctx *gin.Context) {

}

func get(ctx *gin.Context) {

}

func getSingle(ctx *gin.Context) {

}

func update(ctx *gin.Context) {

}

func delete(ctx *gin.Context) {

}







func createIndustryHandler(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	new := models.Industry{
		Name: req.Name,
	}
	industry, err := createIndustry(new)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created industry",
		StatusCode: http.StatusOK,
		Data:       industry,
	})
}

func getIndustryHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	industries, total, page, perPage, err := getIndustry(name, pageSize, pageNumber)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched industries",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     industries,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getSingleIndustryHandler(ctx *gin.Context) {
	industryID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	industry, err := getSingleIndustry(models.Industry{ID: industryID})
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched industry",
		StatusCode: http.StatusOK,
		Data:       industry,
	})
}

func updateIndustryHandler(ctx *gin.Context) {
	industryID, err := uuid.Parse(ctx.Param("id"))
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
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	industry, err := updateIndustry(industryID, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated industry",
		StatusCode: http.StatusOK,
		Data:       industry,
	})
}

func deleteIndustryHandler(ctx *gin.Context) {
	industryID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = deleteSingleIndustry(industryID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted industry",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}





func createSizes(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	new := models.EmployeesSize{
		Name: req.Name,
	}
	industry, err := createEmployeesSize(new)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created employee size",
		StatusCode: http.StatusOK,
		Data:       industry,
	})
}

func getSizes(ctx *gin.Context) {
	name := ctx.Query("name")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	sizes, total, page, perPage, err := getEmployeesSize(name, pageSize, pageNumber)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched employee sizes",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     sizes,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getSingleSizes(ctx *gin.Context) {
	sizeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	size, err := getSingleEmployeesSize(models.EmployeesSize{ID: sizeID})
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched employee size",
		StatusCode: http.StatusOK,
		Data:       size,
	})
}

func updateSizes(ctx *gin.Context) {
	sizeID, err := uuid.Parse(ctx.Param("id"))
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
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	size, err := updateEmployeesSize(sizeID, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated employee size",
		StatusCode: http.StatusOK,
		Data:       size,
	})
}

func deleteSizes(ctx *gin.Context) {
	sizeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = deleteSingle(sizeID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted size",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}