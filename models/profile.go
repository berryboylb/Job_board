package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	// "log"
	// "os"
	"time"
)

type Profile struct {
	gorm.Model
	ID                              uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID                          uuid.UUID              `gorm:"type:uuid;uniqueIndex"` // Unique index enforces one wallet per user
	User                            User                   `gorm:"foreignKey:UserID"`
	Bio                             string                 `gorm:"type:LONGTEXT;not null"`
	Educations                      []Education            `gorm:"foreignKey:ProfileID"`
	EducationsInternShipExperiences []InternShipExperience `gorm:"foreignKey:ProfileID"`
	ProjectsExperiences             []InternShipExperience `gorm:"foreignKey:ProfileID"`
	WorkSamples                     []WorkSample           `gorm:"foreignKey:ProfileID"`
	Awards                          []Award                `gorm:"foreignKey:ProfileID"`
	ProfileLanguages                []ProfileLanguage      `gorm:"foreignKey:ProfileID"`
	SocialMediaAccounts             []SocialMediaAccount   `gorm:"foreignKey:ProfileID"`
}

type Degree struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex"`
}

type Education struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID       uuid.UUID `gorm:"type:uuid;not null"`
	InstitutionName string    `gorm:"type:varchar(100);not null"`
	FieldOFStudy    string    `gorm:"type:varchar(250);not null"`
	DegreeID        uuid.UUID `gorm:"type:uuid;not null"`
	Degree          Degree    `gorm:"foreignKey:DegreeID"`
	GraduationYear  int       `gorm:"not null"`
	StartDate       time.Time
	EndDate         *time.Time
	IsCurrent       *bool `gorm:"default:false"`
}

type InternShipExperience struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID   uuid.UUID `gorm:"type:uuid;not null"`
	CompanyName string    `gorm:"type:varchar(250);not null"`
	Title       string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:LONGTEXT;not null"`
	StartDate   time.Time
	EndDate     *time.Time
	IsCurrent   *bool `gorm:"default:false"`
}

type ProjectsExperience struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID   uuid.UUID `gorm:"type:uuid;not null"`
	ProjectName string    `gorm:"type:varchar(250);not null"`
	Title       string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:LONGTEXT;not null"`
	Link        string    `gorm:"type:varchar(512);not null" ` // Adjusted to use varchar and added not null
	StartDate   time.Time
	EndDate     *time.Time // Made nullable
}

type WorkSample struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID   uuid.UUID `gorm:"type:uuid;not null"`
	Attachment  string    `gorm:"type:varchar(512);not null" `
	Link        string    `gorm:"type:varchar(512);not null" `
	Description string    `gorm:"type:LONGTEXT;not null"`
}

type Award struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID   uuid.UUID `gorm:"type:uuid;not null"`
	Title       string    `gorm:"type:varchar(100);not null"`
	Year        int       `gorm:"not null"`
	Description string    `gorm:"type:LONGTEXT;not null"`
}

type Language struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex"`
}

type LanguageProficiency struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex"`
}

type ProfileLanguage struct {
	// gorm.Model            `gorm:"uniqueIndex:idx_profile_language_proficiency,profile_id,language_proficiency_id"`
	gorm.Model
	ID                    uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID             uuid.UUID           `gorm:"type:uuid;not null"`
	Name                  string              `gorm:"type:varchar(100);not null"`
	LanguageID            uuid.UUID           `gorm:"type:uuid;not null"`
	Language              Language `gorm:"foreignKey:LanguageID"`
	LanguageProficiencyID uuid.UUID           `gorm:"type:uuid;not null"`
	LanguageProficiency   LanguageProficiency `gorm:"foreignKey:LanguageProficiencyID"`
}

type SocialMedia struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex"`
}

type SocialMediaAccount struct {
	gorm.Model
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfileID     uuid.UUID   `gorm:"type:uuid;not null"`
	Link          string      `gorm:"type:varchar(512);not null" `
	SocialMediaID uuid.UUID   `gorm:"type:uuid;not null"`
	SocialMedia   SocialMedia `gorm:"foreignKey:SocialMediaID"`
}
