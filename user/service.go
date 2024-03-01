package user

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

func GetSingleUser(filter models.User) (*models.User, error) {
	var user models.User
	if err := database.Where(&filter).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsers(filter UserDetails, pageNumber string, pageSize string) ([]models.User, int64, int, int, error) {
	perPage := 15
	page := 1

	if pageSize != "" {
		if perPageNum, err := strconv.Atoi(pageSize); err == nil {
			perPage = perPageNum
		}
	}

	if pageNumber != "" {
		if num, err := strconv.Atoi(pageNumber); err == nil {
			page = num
		}
	}

	offset := (page - 1) * perPage

	var users []models.User
	var total int64

	db := database.Model(&models.User{})

	if filter.Name != "" {
		db = db.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Email != "" {
		db = db.Where("email = ?", filter.Email)
	}

	if filter.Picture != "" {
		db = db.Where("picture LIKE ?", "%"+filter.Picture+"%")
	}

	if filter.MobileNumber != "" {
		db = db.Where("mobile_number LIKE ?", "%"+filter.MobileNumber+"%")
	}

	if filter.RoleName != "" {
		db = db.Where("role_name = ?", filter.Email)
	}

	if filter.ProviderID != "" {
		db = db.Where("provider_id = ?", filter.ProviderID)
	}

	if filter.SubscriberID != "" {
		db = db.Where("subscriber_id = ?", filter.SubscriberID)
	}

	if err := db.Count(&total).Error; err != nil {
		log.Println("Error counting books:", err)
		return nil, 0, 0, 0, err
	}

	if err := db.Limit(perPage).Offset(offset).Find(&users).Error; err != nil {
		log.Println("Error finding books:", err)
		return nil, 0, 0, 0, err
	}

	return users, total, page, perPage, nil
}

func DeleteSingleUser(userID uuid.UUID) error {
	err := database.Delete(&models.User{}, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user already deleted")
		}
		return err
	}
	return nil
}

func UpdateSingleUser(id uuid.UUID, values interface{}) (*models.User, error) {
	err := database.Model(models.User{}).Where("id = ?", id).Updates(values).Error
	if err != nil {
		return nil, err
	}

	// Retrieve the updated user
	var user models.User
	err = database.Preload("Role").Preload("Role.Permissions").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateAdminUser(user models.User) (*models.User, error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Try to find the user
    existingUser := &models.User{}
    result := tx.Preload("Profile").Where("mobile_number = ? OR email = ?", user.MobileNumber, user.Email).First(existingUser)
    if result.Error != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error fetching user: %w", result.Error)
    }

    if result.RowsAffected > 0 {
        tx.Rollback()
        return nil, fmt.Errorf("user with the same email or mobile number already exists")
    }

    if err := tx.Create(&user).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error creating user: %w", err)
    }

    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    return &user, nil
}

