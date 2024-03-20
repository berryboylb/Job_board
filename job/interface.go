package job

import (
	"github.com/google/uuid"
	"job_board/models"

	"time"
)

type Request struct {
	Name string `json:"name" binding:"required"`
}

type JobRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CountryID   uuid.UUID `json:"country_id" binding:"required"`
	Salary      float64   `json:"country_id" binding:"required"`
	JobTypeID   uuid.UUID `json:"job_type_id" binding:"required"`
	LevelID     uuid.UUID `json:"level_id" binding:"required"`
	Skills      []string  `json:"skills" binding:"required"`
	CompanyID   uuid.UUID `json:"company_id" binding:"required"`
}

type ApplicationRequest struct {
	JobID uuid.UUID `json:"job_id" binding:"required"`
}

type UpdateApplicationRequest struct {
	Status      string`json:"status" binding:"required"`
}


type SearchApplication struct {
	JobID       uuid.UUID     `json:"job_id" binding:"required"`
	ApplicantID uuid.UUID     `json:"applicant_id" binding:"required"`
	Status      models.Status `json:"status" binding:"required"`
	AppliedAt   time.Time     `json:"time" binding:"required"`
}

