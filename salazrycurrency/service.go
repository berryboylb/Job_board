package salazrycurrency

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"job_board/db"
	"job_board/models"
)

var database *gorm.DB

func init() {
	database = db.GetDB()
}

func createSalary(salaryCurrency models.SalaryCurrency) (*models.SalaryCurrency, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&salaryCurrency).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("salaryCurrency with the same name already exists")
		}
		return nil, fmt.Errorf("error creating a new SalaryCurrency: %v", err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &salaryCurrency, nil
}

func getSalary(name string, pageSize string, pageNumber string) ([]models.SalaryCurrency, int64, int, int, error) {
	// Set default values for page size and page number
	perPage := 15
	page := 1

	// Parse page size and page number if provided
	if pageSize != "" {
		if perPageNum, err := strconv.Atoi(pageSize); err == nil {
			perPage = perPageNum
		}
	}
	if pageNumber != "" {
		if pageNum, err := strconv.Atoi(pageNumber); err == nil {
			page = pageNum
		}
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Initialize database model with filtering conditions
	db := database.Model(&models.SalaryCurrency{})

	// Add WHERE clauses only if the corresponding filter fields are not empty or zero
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	// Count total number of profiles
	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Println("Error counting profiles:", err)
		return nil, 0, 0, 0, err
	}

	// Retrieve profiles with preloaded associations
	var profiles []models.SalaryCurrency
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		log.Println("Error finding SalaryCurrency:", err)
		return nil, 0, 0, 0, err
	}

	return profiles, total, page, perPage, nil
}

func getSingleSalary(ID uuid.UUID) (*models.SalaryCurrency, error) {
	SalaryCurrency := models.SalaryCurrency{}
	if err := database.
		First(&SalaryCurrency, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &SalaryCurrency, nil
}

func updateSalary(ID uuid.UUID, name string) (*models.SalaryCurrency, error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
   var existingRecord models.SalaryCurrency
	if err := tx.First(&existingRecord, "id = ?", ID).Error; err != nil {
		return nil, err // Record not found or other database error
	}

	// Update the record with the provided updates
	if err := tx.Model(&existingRecord).Update("name", name).Error; err != nil {
		return nil, err // Error updating the record
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &existingRecord, nil
}

func deleteSingle(SalaryCurrencyID uuid.UUID) error {
	result := database.Delete(&models.SalaryCurrency{}, SalaryCurrencyID)
	if result.RowsAffected == 0 {
		return errors.New("salary already deleted")
	}
	return result.Error
}
