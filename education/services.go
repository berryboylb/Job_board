package education

import (
	// "errors"
	"fmt"
	// "log"
	// "strconv"

	// "github.com/google/uuid"
	"gorm.io/gorm"

	"job_board/db"
	"job_board/models"
)

var database *gorm.DB

func init() {
	database = db.GetDB()
}

func createEducation(education models.Education) (*models.Education, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&education).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating a new education: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &education, nil
}
