package profile

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

func createProfile(userID uuid.UUID, profile models.Profile) (*models.Profile, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			// Check if the profile already exists (including soft-deleted profiles)
			var existingProfile models.Profile
			if tx.Unscoped().Where("user_id = ?", userID).First(&existingProfile).Error == nil {
				if existingProfile.DeletedAt.Valid {
					// Profile exists but is soft-deleted
					return nil, errors.New("your profile is currently inactive. Please contact support to restore your profile")
				}
				// Profile already exists
				return nil, errors.New("a profile already exists for your account. Please contact support for assistance")
			}
			// Unexpected duplicate error
			return nil, fmt.Errorf("error creating a new profile: %w", err)
		}
		// Other error occurred
		return nil, fmt.Errorf("error creating a new profile: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &profile, nil
}

func getProfiles(filter ProfileDto, pageSize string, pageNumber string) ([]models.Profile, int64, int, int, error) {
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
	db := database.Model(&models.Profile{})

	// Add WHERE clauses only if the corresponding filter fields are not empty or zero
	if filter.Bio != "" {
		db = db.Where("bio LIKE ?", "%"+filter.Bio+"%")
	}
	if filter.Resume != "" {
		db = db.Where("resume LIKE ?", "%"+filter.Resume+"%")
	}
	if filter.GenderID != uuid.Nil {
		db = db.Where("gender_id = ?", filter.GenderID)
	}
	if filter.CurrentSalary != 0.0 {
		db = db.Where("current_salary = ?", filter.CurrentSalary)
	}
	if filter.CurrentSalaryCurrencyID != uuid.Nil {
		db = db.Where("current_salary_currency_id = ?", filter.CurrentSalaryCurrencyID)
	}
	if filter.ExpectedSalary != 0.0 {
		db = db.Where("expected_salary = ?", filter.ExpectedSalary)
	}
	if filter.ExpectedSalaryCurrencyID != uuid.Nil {
		db = db.Where("expected_salary_currency_id = ?", filter.ExpectedSalaryCurrencyID)
	}

	// Count total number of profiles
	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Println("Error counting profiles:", err)
		return nil, 0, 0, 0, err
	}

	// Retrieve profiles with preloaded associations
	var profiles []models.Profile
	if err := db.
		// Preload("Educations").
		// Preload("InternShipExperiences").
		// Preload("ProjectsExperiences").
		// Preload("WorkSamples").
		// Preload("Awards").
		// Preload("ProfileLanguages").
		// Preload("SocialMediaAccounts").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		log.Println("Error finding profiles:", err)
		return nil, 0, 0, 0, err
	}

	return profiles, total, page, perPage, nil
}

func getProfile(search models.Profile) (*models.Profile, error) {
	profile := models.Profile{}
	if err := database.
		Preload("Educations").
		Preload("InternShipExperiences").
		Preload("ProjectsExperiences").
		Preload("WorkSamples").
		Preload("Awards").
		Preload("ProfileLanguages").
		Preload("SocialMediaAccounts").
		Where(&search).
		First(&profile).Error; 
		err != nil {
		return nil, err
	}

	return &profile, nil
}



func updateProfile(search models.Profile, profile map[string]interface{}) (*models.Profile, error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.
        Model(&models.Profile{}).
        Where(&search).
        Updates(profile).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error updating profile: %w", err)
    }

    if err := tx.Commit().Error; err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    updatedProfile := models.Profile{}
    if err := tx.
        Where(&search).
		Preload("Educations").
        Preload("InternShipExperiences").
        Preload("ProjectsExperiences").
        Preload("WorkSamples").
        Preload("Awards").
        Preload("ProfileLanguages").
        Preload("SocialMediaAccounts").
        First(&updatedProfile).Error; err != nil {
        return nil, fmt.Errorf("error fetching updated profile: %w", err)
    }

    return &updatedProfile, nil
}



func deleteSingleProfile(search models.Profile) error {
	result := database.Delete(&search)
	if result.RowsAffected == 0 {
		return errors.New("user already deleted")
	}
	return result.Error
}