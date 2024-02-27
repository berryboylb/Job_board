package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Industry struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(250);not null;unique"`
}

type EmployeesSize struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex"`
}

type Company struct {
	gorm.Model
	ID              uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name            string        `gorm:"type:varchar(100);not null"`
	Description     string        `gorm:"type:text;not null"`
	Website         string        `gorm:"type:varchar(512)"`
	IndustryID      uuid.UUID     `gorm:"type:uuid;not null"`
	Industry        Industry      `gorm:"foreignKey:IndustryID"`
	Established     time.Time     // Year the company was established
	Location        string        `gorm:"type:varchar(100)"` // Location of the company
	EmployeesSizeID uuid.UUID     `gorm:"type:uuid;not null"`
	EmployeesSize   EmployeesSize `gorm:"foreignKey: EmployeesSizeID"`
	Logo            string        `gorm:"type:varchar(512);default:'https://via.placeholder.com/200x200'"`
	UserID          uuid.UUID     `gorm:"type:uuid;not null"` // Removed uniqueIndex
	User            User          `gorm:"foreignKey:UserID"`
}

type JobType struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex"`
}

type Level struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex"`
}

type Job struct {
	gorm.Model
	ID              uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title           string           `gorm:"type:varchar(100);not null"`
	Description     string           `gorm:"type:text;not null"`
	Location        string           `gorm:"type:varchar(250);not null"`
	Salary          float64          `gorm:"type:decimal(10,2);default:0.0"`
	JobTypeID       uuid.UUID        `gorm:"type:uuid;not null"`
	JobType         JobType          `gorm:"foreignKey:JobTypeID"`
	LevelID         uuid.UUID        `gorm:"type:uuid;not null"`
	Level           Level            `gorm:"foreignKey:LevelID"`
	Skills          pq.StringArray   `json:"skills" gorm:"type:text[]; not null"`
	CompanyID       uuid.UUID        `gorm:"type:uuid;not null"`
	Company         Company          `gorm:"foreignKey: CompanyID"`
	JobApplications []JobApplication `gorm:"foreignKey:JobID"`
}

type Status string

const (
	Pending Status = "pending"
	Success Status = "success"
	Failed  Status = "failed"
)

/*
you need to run this cript to ensure a user can only apply once to a job

ALTER TABLE job_applications
ADD CONSTRAINT unique_applicant_job UNIQUE (applicant_id, job_id);
*/
type JobApplication struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	JobID       uuid.UUID `gorm:"type:uuid;not null;index:idx_job_applications"`
	Job         Job       `gorm:"foreignKey:JobID"`
	ApplicantID uuid.UUID `gorm:"type:uuid;not null;index:idx_job_applications"`
	Applicant   User      `gorm:"foreignKey:ApplicantID"`
	Status      Status    `gorm:"type:varchar(50);default:'pending'"`
	AppliedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
