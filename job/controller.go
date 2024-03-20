package job

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"net/http"
	"strconv"
	"strings"
	"time"

	"job_board/helpers"
	"job_board/models"
)

func create(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	var req JobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	newJob := models.Job{
		UserID:      user.ID,
		Title:       req.Title,
		Description: req.Description,
		CountryID:   req.CountryID,
		Salary:      req.Salary,
		JobTypeID:   req.JobTypeID,
		LevelID:     req.LevelID,
		Skills:      pq.StringArray(req.Skills),
		CompanyID:   req.CompanyID,
	}

	resp, err := createJob(newJob)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully created job",
		StatusCode: http.StatusOK,
		Data:       resp,
	})

	/* job alerts for everyone in this location */
}

func get(ctx *gin.Context) {
	var (
		countryID uuid.UUID
		jobTypeID uuid.UUID
		levelID   uuid.UUID
		companyID uuid.UUID
		skills    []string
		salary    float64
		err       error
		errsArr  []string
	)

	if id := ctx.Query("country_id"); id != "" {
		if countryID, err = uuid.Parse(id); err != nil {
			errsArr = append(errsArr, err.Error())
		}
	}

	if id := ctx.Query("job_type_id"); id != "" {
		if jobTypeID, err = uuid.Parse(id); err != nil {
			errsArr = append(errsArr, err.Error())
		}
	}

	if id := ctx.Query("level_id"); id != "" {
		if levelID, err = uuid.Parse(id); err != nil {
			errsArr = append(errsArr, err.Error())
		}
	}

	if id := ctx.Query("company_id"); id != "" {
		if companyID, err = uuid.Parse(id); err != nil {
			errsArr = append(errsArr, err.Error())
		}
	}
	if skill := ctx.Query("skill"); skill != "" {
		skills = append(skills, skill)
	}
	if val := ctx.Query("salary"); val != "" {
		if salary, err = strconv.ParseFloat(val, 10); err != nil {
			errsArr = append(errsArr, err.Error())
		}
	}
	if len(errsArr) > 0 {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    strings.Join(errsArr, ","),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	filter := JobRequest{
		CompanyID:   companyID,
		JobTypeID:   jobTypeID,
		LevelID:     levelID,
		CountryID:   countryID,
		Skills:      skills,
		Salary:      salary,
		Title:       ctx.Query("title"),
		Description: ctx.Query("description"),
	}

	resp, total, page, perPage, err := getJob(filter, ctx.Query("page_size"), ctx.Query("page_number"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched jobs",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     resp,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})

}

func getSingle(ctx *gin.Context) {
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

	resp, err := getSingleJob(ID, *user)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully fetched job",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func update(ctx *gin.Context) {
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
	var req JobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	resp, err := updateJob(ID, *user, req)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully updated job",
		StatusCode: http.StatusOK,
		Data:       resp,
	})

	/* if the country is being change and the job is less than 24 hours send alert*/
}

func delete(ctx *gin.Context) {
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

	err = deleteSingleJob(ID, *user)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully deleted job",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func createLevel(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	new := models.Level{
		Name: req.Name,
	}
	level, err := createLevelFunc(new)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created level",
		StatusCode: http.StatusOK,
		Data:       level,
	})
}

func getLevel(ctx *gin.Context) {
	name := ctx.Query("name")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	levels, total, page, perPage, err := getLevelFunc(name, pageSize, pageNumber)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched levels",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     levels,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getSingleLevel(ctx *gin.Context) {
	levelID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	level, err := getSingleLevelFunc(models.Level{ID: levelID})
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched level",
		StatusCode: http.StatusOK,
		Data:       level,
	})
}

func updateLevel(ctx *gin.Context) {
	levelID, err := uuid.Parse(ctx.Param("id"))
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
	level, err := updateLevelFunc(levelID, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated level",
		StatusCode: http.StatusOK,
		Data:       level,
	})
}

func deleteLevel(ctx *gin.Context) {
	levelID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = deleteSingle(levelID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted level",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func createType(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	new := models.JobType{
		Name: req.Name,
	}
	jobtype, err := createJobType(new)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully created job type",
		StatusCode: http.StatusOK,
		Data:       jobtype,
	})
}

func getType(ctx *gin.Context) {
	name := ctx.Query("name")
	pageSize := ctx.Query("page_size")
	pageNumber := ctx.Query("page_number")
	types, total, page, perPage, err := getJobType(name, pageSize, pageNumber)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched types",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     types,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getSingleType(ctx *gin.Context) {
	typeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	types, err := getSingleJobType(models.JobType{ID: typeID})
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched job type",
		StatusCode: http.StatusOK,
		Data:       types,
	})
}

func updateType(ctx *gin.Context) {
	typeID, err := uuid.Parse(ctx.Param("id"))
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
	types, err := updateLevelFunc(typeID, req.Name)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully updated job type",
		StatusCode: http.StatusOK,
		Data:       types,
	})
}

func deleteType(ctx *gin.Context) {
	typeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	err = deleteSingleJobType(typeID)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully deleted job type",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

func createApplication(ctx *gin.Context) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	var req ApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	newProject := models.JobApplication{
		JobID: req.JobID,
		ApplicantID: user.ID,
		Status: models.Pending,
		AppliedAt: time.Now() ,
	}

	resp, err := createJobApplication(newProject, *user)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully created application",
		StatusCode: http.StatusOK,
		Data:       resp,
	})

	// notify poster with novu
}

func getApplication(ctx *gin.Context) {
	var (
		jobID  uuid.UUID
		applicantID  uuid.UUID
		status      models.Status         
		appliedAt   time.Time 
		err       error
	)
	if date := ctx.Query("applied_at"); date != "" {
		if appliedAt, err = time.Parse("2006-01-02", date); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if id := ctx.Query("job_id"); id != "" {
		if jobID, err = uuid.Parse(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

		if id := ctx.Query("applicant_id"); id != "" {
		if applicantID, err = uuid.Parse(id); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	if sta := ctx.Query("status"); sta != "" {
		if status, err = models.ParseStatus(sta); err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
	}

	filter := SearchApplication{
		JobID: jobID,
		ApplicantID: applicantID ,
		Status: status,
		AppliedAt: appliedAt ,
	}

	resp, total, page, perPage, err := getJobApplication(filter, ctx.Query("page_size"), ctx.Query("page_number"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched application",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     resp,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getPosterJobApplication(ctx *gin.Context) {
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

	resp, total, page, perPage, err := getApplications(ID, *user, ctx.Query("page_size"), ctx.Query("page_number"))
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}
	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "successfully fetched applications",
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"data":     resp,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		},
	})
}

func getSingleApplication(ctx *gin.Context) {
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

	resp, err := getSingleJobApplication(ID, *user)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully fetched applications",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func updateApplication(ctx *gin.Context) {
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
	var req UpdateApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	 status, err := models.ParseStatus(req.Status);
	 if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	if err = models.CheckStatus(status); err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	resp, err := updateJobApplication(ID, *user, status)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully updated status",
		StatusCode: http.StatusOK,
		Data:       resp,
	})
}

func deleteApplication(ctx *gin.Context) {
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

	err = deleteSingleJobApplication(ID, *user)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	helpers.CreateResponse(ctx, helpers.Response{
		Message:    "Successfully deleted application",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}
