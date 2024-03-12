package language

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"

	"job_board/helpers"
	"job_board/models"
)

/*language*/

func CreateLanguage(ctx *gin.Context) {
	var req Language
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	new := models.Language{
		Name: req.Name,
	}
	language, err := createLanguage(new)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created language",
		StatusCode: http.StatusOK,
		Data:       language,
	})
}

func GetLanguage(ctx *gin.Context) {
	name := ctx.Query("name")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	language, total, page, perPage, err := getLanguage(name, pageSize, pageNumber)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched language",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     language,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func GetSingleLanguage(ctx *gin.Context) {
	languageID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	language, err := getSingleLanguage(models.Language{ID: languageID})
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched language",
		StatusCode: http.StatusOK,
		Data:       language,
	})
}

func UpdateLanguage(ctx *gin.Context) {
	languageID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	var req Language
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	language, err := updateLanguage(languageID, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated language",
		StatusCode: http.StatusOK,
		Data:       language,
	})
}

func DeleteLanguage(ctx *gin.Context) {
	languageID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = deleteLanguage(languageID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted language",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

/*language segment */

/*language*/

func CreateLanguageProficiency(ctx *gin.Context) {
	var req Language
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	new := models.LanguageProficiency{
		Name: req.Name,
	}
	proficiency, err := createLanguageProficiency(new)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created proficiency",
		StatusCode: http.StatusOK,
		Data:       proficiency,
	})
}

func GetLanguageProficiency(ctx *gin.Context) {
	name := ctx.Query("name")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	proficiency, total, page, perPage, err := getLanguageProficiency(name, pageSize, pageNumber)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched proficiency",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     proficiency,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func GetSingleLanguageProficiency(ctx *gin.Context) {
	languageProficiencyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	languageProficiency, err := getSingleLanguageProficiency(models.LanguageProficiency{ID: languageProficiencyID})
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched language proficiency",
		StatusCode: http.StatusOK,
		Data:       languageProficiency,
	})
}

func UpdateLanguageProficiency(ctx *gin.Context) {
	languageProficiencyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	var req Language
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	languageProficiency, err := updateLanguage(languageProficiencyID, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated language proficiency",
		StatusCode: http.StatusOK,
		Data:       languageProficiency,
	})
}

func DeleteLanguageProficiency(ctx *gin.Context) {
	languageProficiencyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = deleteSingle(languageProficiencyID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted language Proficiency",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

/*language segment */

/* profile language segment  starts*/

func CreateProfileLanguage(ctx *gin.Context) {
	
}

func GetProfileLanguage(ctx *gin.Context) {
	
}

func GetSingleProfileLanguage(ctx *gin.Context) {

}

func UpdateProfileLanguage(ctx *gin.Context) {

}

func DeleteProfileLanguage(ctx *gin.Context) {

}

/* ProfileLanguage segment ends*/
