package models

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "github.com/lib/pq"
	"gorm.io/gorm"

	"errors"
	"time"

	"job_board/helpers"
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
	ID                uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name              string           `json:"name"`
	Email             string           `json:"email"`
	Picture           string           `json:"picture"`
	RoleName          RoleAllowed      `json:"role" sql:"type:role_name"`
	ProviderID        string           `gorm:"uniqueIndex" json:"provider_id"`
	Profile           *Profile         `gorm:"foreignKey:UserID" json:"profile"`
	MobileNumber      *string          `gorm:"type:varchar(25);default:null" json:"mobile_number"`
	SubscriberID      string           `gorm:"default:null" json:"subscriber_id"`
	JobApplications   []JobApplication `gorm:"foreignKey:ApplicantID" json:"job_applications"`
	Jobs              []Job            `gorm:"foreignKey:UserID" json:"jobs"`
	Companies         []Company        `gorm:"foreignKey:UserID" json:"companies"`
	VerificationToken string           `json:"verification_token"`
	ExpiresAt         time.Time        `json:"expires_at"`
	Password          string           `gorm:"default:null" json:"-"`
	CountryID         uuid.UUID        `gorm:"type:uuid;"`
	Country           Country          `gorm:"foreignKey:CountryID"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.RoleName == SuperAdminRole {
		// Check if the password is not empty
		if u.Password != "" {
			// Hash Password
			u.Password, err = helpers.HashPassword(u.Password, 14)
			if err != nil {
				return err
			}
		}
	}

	if u.RoleName == AdminRole {
		// Check if the password is not empty
		if u.Password != "" {
			// Hash Password
			u.Password, err = helpers.HashPassword(u.Password, 10)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetUserFromContext(ctx *gin.Context) (*User, error) {
	value, exists := ctx.Get("user")
	if !exists {
		return nil, errors.New("user not found in session")
	}

	user, ok := value.(User)
	if !ok {
		return nil, errors.New("Mismatching types")
	}
	return &user, nil
}
