package company

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

func createIndustry(Industry  models.Industry ) (*models.Industry , error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&Industry ).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("Industry  with the same name already exists")
		}
		return nil, fmt.Errorf("error creating a new Industry : %v", err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &Industry , nil
}

func getIndustry (name string, pageSize string, pageNumber string) ([]models.Industry , int64, int, int, error) {
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
	db := database.Model(&models.Industry {})

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
	var profiles []models.Industry 
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		log.Println("Error finding Industry :", err)
		return nil, 0, 0, 0, err
	}

	return profiles, total, page, perPage, nil
}

func getSingleIndustry (search models.Industry ) (*models.Industry , error) {
	Industry  := models.Industry {}
	if err := database.
		Where(&search).
		First(&Industry ).Error; err != nil {
		return nil, err
	}
	return &Industry , nil
}

func updateIndustry (IndustryID uuid.UUID, name string) (*models.Industry , error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    result := tx.Model(&models.Industry {}).Where("id = ?", IndustryID).Update("name", name)
    if result.Error != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error updating Industry : %w", result.Error)
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Return the updated Industry 
    updatedIndustry  := &models.Industry {ID: IndustryID, Name: name}
    return updatedIndustry , nil
}

func deleteSingleIndustry(IndustryID uuid.UUID) error {
	result := database.Delete(&models.Industry{}, IndustryID)
	if result.RowsAffected == 0 {
		return errors.New("Industry  already deleted")
	}
	return result.Error
}



func createEmployeesSize(EmployeesSize models.EmployeesSize) (*models.EmployeesSize, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&EmployeesSize).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("EmployeesSize with the same name already exists")
		}
		return nil, fmt.Errorf("error creating a new EmployeesSize: %v", err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &EmployeesSize, nil
}

func getEmployeesSize(name string, pageSize string, pageNumber string) ([]models.EmployeesSize, int64, int, int, error) {
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
	db := database.Model(&models.EmployeesSize{})

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
	var profiles []models.EmployeesSize
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		log.Println("Error finding EmployeesSize:", err)
		return nil, 0, 0, 0, err
	}

	return profiles, total, page, perPage, nil
}

func getSingleEmployeesSize(search models.EmployeesSize) (*models.EmployeesSize, error) {
	EmployeesSize := models.EmployeesSize{}
	if err := database.
		Where(&search).
		First(&EmployeesSize).Error; err != nil {
		return nil, err
	}
	return &EmployeesSize, nil
}

func updateEmployeesSize(EmployeesSizeID uuid.UUID, name string) (*models.EmployeesSize, error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    result := tx.Model(&models.EmployeesSize{}).Where("id = ?", EmployeesSizeID).Update("name", name)
    if result.Error != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error updating EmployeesSize: %w", result.Error)
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Return the updated EmployeesSize
    updatedEmployeesSize := &models.EmployeesSize{ID: EmployeesSizeID, Name: name}
    return updatedEmployeesSize, nil
}

func deleteSingle(EmployeesSizeID uuid.UUID) error {
	result := database.Delete(&models.EmployeesSize{}, EmployeesSizeID)
	if result.RowsAffected == 0 {
		return errors.New("EmployeesSize already deleted")
	}
	return result.Error
}