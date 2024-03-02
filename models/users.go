package models

import (
	"github.com/google/uuid"
	// "github.com/lib/pq"
	"gorm.io/gorm"

	"time"
)

type RoleAllowed string

const (
	SuperAdminRole RoleAllowed = "super-admin"
	AdminRole      RoleAllowed = "admin"
	PosterRole     RoleAllowed = "poster"
	UserRole       RoleAllowed = "user"
)

// user struct
type User struct {
	gorm.Model
	ID                uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name              string           // Assuming no special column name needed
	Email             string           `json:"email"`
	Picture           string           `json:"picture"`
	RoleName          RoleAllowed      `json:"role" sql:"type:role_name"`
	ProviderID        string           `gorm:"uniqueIndex"`
	Profile           *Profile         `gorm:"foreignKey:UserID"`
	MobileNumber      *string          `gorm:"type:varchar(25);default:null"` // Assuming unique mobile numbers
	SubscriberID      string           `gorm:"default:null"`
	JobApplications   []JobApplication `gorm:"foreignKey:ApplicantID"`
	Companies         []Company        `gorm:"foreignKey:UserID"`
	VerificationToken string
	ExpiresAt         time.Time
	Password          string `gorm:"default:null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	if u.RoleName == SuperAdminRole || u.RoleName == AdminRole {
		//HASH Password
	}
	return
}
