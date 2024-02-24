package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	// "log"
	// "os"
	// "time"
)

type JobType struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `gorm:"type:varchar(250);not null; uniqueIndex"`
}

type Jobs struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title          string    `gorm:"type:varchar(100);not null"`
	Description    string    `gorm:"type:LONGTEXT;not null"`
	Location       string    `gorm:"type:varchar(250);not null"`
	Salary         float64   `gorm:"type:decimal(10,2);default:0.0"`
	JobTypeID      uuid.UUID `gorm:"type:uuid;not null"`
	JobType        JobType   `gorm:"foreignKey:JobTypeID"`
}
