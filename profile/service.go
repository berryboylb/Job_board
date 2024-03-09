package profile

import (
	"errors"
	"fmt"

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

