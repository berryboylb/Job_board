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

func getSingleLanguage(ID uuid.UUID) (*models.Language, error) {
	Language := models.Language{}
	if err := database.
		First(&Language, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &Language, nil
}

func updateLanguage(ID uuid.UUID, name string) (*models.Language, error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    var existingRecord models.Language
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

func getSingleLanguageProficiency(ID uuid.UUID) (*models.LanguageProficiency, error) {
	LanguageProficiency := models.LanguageProficiency{}
	if err := database.
		First(&LanguageProficiency, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &LanguageProficiency, nil
}

func updateLanguageProficiency(ID uuid.UUID, name string) (*models.LanguageProficiency, error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
   var existingRecord models.LanguageProficiency
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

func deleteSingle(LanguageProficiencyID uuid.UUID) error {
	result := database.Delete(&models.LanguageProficiency{}, LanguageProficiencyID)
	if result.RowsAffected == 0 {
		return errors.New("LanguageProficiency already deleted")
	}
	return result.Error
}

func checkProfile(user models.User) bool {
	if user.Profile == nil {
		return true
	}
	return false
}

func createProficiency(Proficiency models.ProfileLanguage, user models.User) (*models.ProfileLanguage, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if profile := checkProfile(user); profile {
		return nil, fmt.Errorf("you don't have a profile")
	}

	if err := tx.Create(&Proficiency).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating a new Proficiency experience: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &Proficiency, nil
}

func getProficiency(filter Search, pageSize string, pageNumber string) ([]models.ProfileLanguage, int64, int, int, error) {
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
	db := database.Model(&models.ProfileLanguage{})

	// Add WHERE clauses only if the corresponding filter fields are not empty or zero
	if filter.Name != "" {
		db = db.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.LanguageID != uuid.Nil {
		db = db.Where("language_id = ?", filter.LanguageID)
	}
	if filter.LanguageProficiencyID != uuid.Nil {
		db = db.Where("language_proficiency_id = ?", filter.LanguageProficiencyID)
	}

	// Count total number of profiles
	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Println("Error counting profiles:", err)
		return nil, 0, 0, 0, err
	}

	// Retrieve profiles with preloaded associations
	var data []models.ProfileLanguage
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&data).Error; err != nil {
		log.Println("Error finding Proficiency:", err)
		return nil, 0, 0, 0, err
	}

	return data, total, page, perPage, nil
}

func getSingleProficiency(ID uuid.UUID, user models.User) (*models.ProfileLanguage, error) {
	if profile := checkProfile(user); profile {
		return nil, fmt.Errorf("you don't have a profile")
	}
	var record models.ProfileLanguage
	if err := database.
		Preload("Language").
		Preload("LanguageProficiency").
		First(&record, "id = ?", ID).Error; err != nil {
		return nil, err
	}

	// Check if the user has permission to update the record
	if user.RoleName == models.UserRole && record.ProfileID != user.Profile.ID {
		return nil, fmt.Errorf("you don't have permission to view this record")
	}

	return &record, nil
}

func updateProficiency(ID uuid.UUID, user models.User, updates Request) (*models.ProfileLanguage, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingRecord models.ProfileLanguage
	if err := tx.First(&existingRecord, "id = ?", ID).Error; err != nil {
		return nil, err // Record not found or other database error
	}

	// Check if the user has permission to update the record
	if user.RoleName == models.UserRole {
		if profile := checkProfile(user); profile {
			return nil, fmt.Errorf("you don't have a profile")
		}
		if existingRecord.ProfileID != user.Profile.ID {
			return nil, fmt.Errorf("you don't have permission to update this record")
		}
	}

	// Update the record with the provided updates
	if err := tx.Model(&existingRecord).Updates(updates).Error; err != nil {
		return nil, err // Error updating the record
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &existingRecord, nil
}

func deleteSingleProficiency(ID uuid.UUID, user models.User) error {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingRecord models.ProfileLanguage
	if err := tx.First(&existingRecord, "id = ?", ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if user.RoleName == models.UserRole {
		if existingRecord.ProfileID != user.Profile.ID {
			return fmt.Errorf("you don't have permission to update this record")
		}
	}

	result := tx.Delete(&models.ProfileLanguage{}, "id = ?", ID)
	if result.Error != nil {
		// Rollback the transaction if an error occurs
		tx.Rollback()
		return result.Error
	}

	// Check if the record was not found
	if result.RowsAffected == 0 {
		// Rollback the transaction and return a custom error
		tx.Rollback()
		return fmt.Errorf("record with ID %s not found", ID)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}