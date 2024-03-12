package language

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

func createLanguage(Language models.Language) (*models.Language, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&Language).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("language with the same name already exists")
		}
		return nil, fmt.Errorf("error creating a new Language: %v", err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &Language, nil
}

func getLanguage(name string, pageSize string, pageNumber string) ([]models.Language, int64, int, int, error) {
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
	db := database.Model(&models.Language{})

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
	var profiles []models.Language
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		log.Println("Error finding Language:", err)
		return nil, 0, 0, 0, err
	}

	return profiles, total, page, perPage, nil
}

func getSingleLanguage(search models.Language) (*models.Language, error) {
	Language := models.Language{}
	if err := database.
		Where(&search).
		First(&Language).Error; err != nil {
		return nil, err
	}
	return &Language, nil
}

func updateLanguage(LanguageID uuid.UUID, name string) (*models.Language, error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    result := tx.Model(&models.Language{}).Where("id = ?", LanguageID).Update("name", name)
    if result.Error != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error updating Language: %w", result.Error)
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Return the updated Language
    updatedLanguage := &models.Language{ID: LanguageID, Name: name}
    return updatedLanguage, nil
}

func deleteLanguage(LanguageID uuid.UUID) error {
	result := database.Delete(&models.Language{}, LanguageID)
	if result.RowsAffected == 0 {
		return errors.New("Language already deleted")
	}
	return result.Error
}



func createLanguageProficiency(LanguageProficiency models.LanguageProficiency) (*models.LanguageProficiency, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&LanguageProficiency).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("LanguageProficiency with the same name already exists")
		}
		return nil, fmt.Errorf("error creating a new LanguageProficiency: %v", err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &LanguageProficiency, nil
}

func getLanguageProficiency(name string, pageSize string, pageNumber string) ([]models.LanguageProficiency, int64, int, int, error) {
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
	db := database.Model(&models.LanguageProficiency{})

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
	var profiles []models.LanguageProficiency
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		log.Println("Error finding LanguageProficiency:", err)
		return nil, 0, 0, 0, err
	}

	return profiles, total, page, perPage, nil
}

func getSingleLanguageProficiency(search models.LanguageProficiency) (*models.LanguageProficiency, error) {
	LanguageProficiency := models.LanguageProficiency{}
	if err := database.
		Where(&search).
		First(&LanguageProficiency).Error; err != nil {
		return nil, err
	}
	return &LanguageProficiency, nil
}

func updateLanguageProficiency(LanguageProficiencyID uuid.UUID, name string) (*models.LanguageProficiency, error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    result := tx.Model(&models.LanguageProficiency{}).Where("id = ?", LanguageProficiencyID).Update("name", name)
    if result.Error != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error updating LanguageProficiency: %w", result.Error)
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Return the updated LanguageProficiency
    updatedLanguageProficiency := &models.LanguageProficiency{ID: LanguageProficiencyID, Name: name}
    return updatedLanguageProficiency, nil
}

func deleteSingle(LanguageProficiencyID uuid.UUID) error {
	result := database.Delete(&models.LanguageProficiency{}, LanguageProficiencyID)
	if result.RowsAffected == 0 {
		return errors.New("LanguageProficiency already deleted")
	}
	return result.Error
}