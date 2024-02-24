package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	// "log"
	// "os"
)

type RoleAllowed string

const (
	AdminRole  RoleAllowed = "admin"
	PosterRole RoleAllowed = "poster"
	UserRole   RoleAllowed = "user"
)

// user struct
type User struct {
	gorm.Model
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name       string         `gorm:"name"`
	Email      string         `gorm:"uniqueIndex"`
	Picture    string         `json:"picture"`
	RoleName   RoleAllowed    `json:"role" sql:"type:role_name"`
	RoleAccess pq.StringArray `json:"role_access" gorm:"type:text[]; not null"`
	Profile    *Profile       `gorm:"foreignKey:UserID"`
}
