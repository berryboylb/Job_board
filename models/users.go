package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	// "log"
	// "os"
	// "time"
)

//user struct
type User struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}






