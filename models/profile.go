package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Profile struct {
	gorm.Model
	ID                       uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID                   uuid.UUID              `gorm:"type:uuid;uniqueIndex"` // Unique index enforces one wallet per user
	User                     User                   `gorm:"foreignKey:UserID"`
	Bio                      string                 `gorm:"type:text;not null"`
	Resume                   string                 `gorm:"type:varchar(512);not null"`
	Educations               []Education            `gorm:"foreignKey:ProfileID"`
	InternShipExperiences    []InternShipExperience `gorm:"foreignKey:ProfileID"`
	ProjectsExperiences      []ProjectsExperience   `gorm:"foreignKey:ProfileID"`
	WorkSamples              []WorkSample           `gorm:"foreignKey:ProfileID"`
	Awards                   []Award                `gorm:"foreignKey:ProfileID"`
	ProfileLanguages         []ProfileLanguage      `gorm:"foreignKey:ProfileID"`
	SocialMediaAccounts      []SocialMediaAccount   `gorm:"foreignKey:ProfileID"`
	GenderID                 uuid.UUID              `gorm:"type:uuid;not null"`
	Gender                   Gender                 `gorm:"foreignKey:GenderID"`
	CurrentSalary            float64                `gorm:"type:decimal(10,2);default:0.0"`
	CurrentSalaryCurrencyID  *uuid.UUID             `gorm:"type:uuid"`
	CurrentSalaryCurrency    SalaryCurrency         `gorm:"foreignKey:CurrentSalaryCurrencyID"`
	ExpectedSalary           float64                `gorm:"type:decimal(10,2);default:0.0"`
	ExpectedSalaryCurrencyID *uuid.UUID             `gorm:"type:uuid"`
	ExpectedSalaryCurrency   SalaryCurrency         `gorm:"foreignKey:ExpectedSalaryCurrencyID"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type SalaryCurrency struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex"json:"name"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type Gender struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex" json:"name"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type Degree struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex" json:"name"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type AcademicRanking struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex" json:"name"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type Education struct {
	gorm.Model
	ID                uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProfileID         uuid.UUID       `gorm:"type:uuid;not null"`
	InstitutionName   string          `gorm:"type:varchar(100);not null"`
	FieldOFStudy      string          `gorm:"type:varchar(250);not null"`
	DegreeID          uuid.UUID       `gorm:"type:uuid;not null"`
	Degree            Degree          `gorm:"foreignKey:DegreeID"`
	AcademicRankingID uuid.UUID       `gorm:"type:uuid;not null"`
	AcademicRanking   AcademicRanking `gorm:"foreignKey:AcademicRankingID"`
	GraduationYear    int             `gorm:"not null"`
	StartDate         time.Time
	EndDate           *time.Time
	IsCurrent         *bool  `gorm:"default:false"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type InternShipExperience struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProfileID   uuid.UUID `gorm:"type:uuid;not null"`
	CompanyName string    `gorm:"type:varchar(250);not null"`
	Title       string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text;not null"`
	StartDate   time.Time
	EndDate     *time.Time
	IsCurrent   *bool `gorm:"default:false"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type ProjectsExperience struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProfileID   uuid.UUID `gorm:"type:uuid;not null"`
	ProjectName string    `gorm:"type:varchar(250);not null"`
	Title       string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text;not null"`
	Link        string    `gorm:"type:varchar(512);not null" ` // Adjusted to use varchar and added not null
	StartDate   time.Time
	EndDate     *time.Time // Made nullable
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type WorkSample struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProfileID   uuid.UUID `gorm:"type:uuid;not null"`
	Attachment  string    `gorm:"type:varchar(512);not null" `
	Link        string    `gorm:"type:varchar(512);not null" `
	Description string    `gorm:"type:text;not null"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type Award struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProfileID   uuid.UUID `gorm:"type:uuid;not null"`
	Title       string    `gorm:"type:varchar(100);not null"`
	Year        int       `gorm:"not null"`
	Description string    `gorm:"type:text;not null"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type Language struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex" json:"name"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type LanguageProficiency struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id" `
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex" json:"name"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type ProfileLanguage struct {
	// gorm.Model            `gorm:"uniqueIndex:idx_profile_language_proficiency,profile_id,language_proficiency_id"`
	gorm.Model
	ID                    uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProfileID             uuid.UUID           `gorm:"type:uuid;not null" json:"name"`
	Name                  string              `gorm:"type:varchar(100);not null"`
	LanguageID            uuid.UUID           `gorm:"type:uuid;not null"`
	Language              Language            `gorm:"foreignKey:LanguageID"`
	LanguageProficiencyID uuid.UUID           `gorm:"type:uuid;not null"`
	LanguageProficiency   LanguageProficiency `gorm:"foreignKey:LanguageProficiencyID"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type SocialMedia struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex" json:"name"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type SocialMediaAccount struct {
	gorm.Model
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProfileID     uuid.UUID   `gorm:"type:uuid;not null"`
	Link          string      `gorm:"type:varchar(512);not null" `
	SocialMediaID uuid.UUID   `gorm:"type:uuid;not null"`
	SocialMedia   SocialMedia `gorm:"foreignKey:SocialMediaID"`
	CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
    DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}
