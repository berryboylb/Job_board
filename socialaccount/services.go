package socialaccount

import (
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

func checkProfile(user models.User) bool {
	if user.Profile == nil {
		return true
	}
	return false
}

func createSocial(Social models.SocialMediaAccount, user models.User) (*models.SocialMediaAccount, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if profile := checkProfile(user); profile {
		return nil, fmt.Errorf("you don't have a profile")
	}

	if err := tx.Create(&Social).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating a new Social experience: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &Social, nil
}

func getSocial(filter Search, pageSize string, pageNumber string) ([]models.SocialMediaAccount, int64, int, int, error) {
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
	db := database.Model(&models.SocialMediaAccount{})

	// Add WHERE clauses only if the corresponding filter fields are not empty or zero
	if filter.Link != "" {
		db = db.Where("link LIKE ?", "%"+filter.Link+"%")
	}
	if filter.SocialMediaID != uuid.Nil {
		db = db.Where("social_media_id = ?", filter.SocialMediaID)
	}

	// Count total number of profiles
	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Println("Error counting profiles:", err)
		return nil, 0, 0, 0, err
	}

	// Retrieve profiles with preloaded associations
	var data []models.SocialMediaAccount
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&data).Error; err != nil {
		log.Println("Error finding Social:", err)
		return nil, 0, 0, 0, err
	}

	return data, total, page, perPage, nil
}

func getSingleSocial(ID uuid.UUID, user models.User) (*models.SocialMediaAccount, error) {
	if profile := checkProfile(user); profile {
		return nil, fmt.Errorf("you don't have a profile")
	}
	var record models.SocialMediaAccount
	if err := database.
		Preload("SocialMediaID").
		First(&record, "id = ?", ID).Error; err != nil {
		return nil, err
	}

	// Check if the user has permission to update the record
	if user.RoleName == models.UserRole && record.ProfileID != user.Profile.ID {
		return nil, fmt.Errorf("you don't have permission to view this record")
	}

	return &record, nil
}

func updateSocial(ID uuid.UUID, user models.User, updates Request) (*models.SocialMediaAccount, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingRecord models.SocialMediaAccount
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

func deleteSingleSocial(ID uuid.UUID, user models.User) error {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingRecord models.SocialMediaAccount
	if err := tx.First(&existingRecord, "id = ?", ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if user.RoleName == models.UserRole {
		if existingRecord.ProfileID != user.Profile.ID {
			return fmt.Errorf("you don't have permission to update this record")
		}
	}

	result := tx.Delete(&models.SocialMediaAccount{}, "id = ?", ID)
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
