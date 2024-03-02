package user

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"job_board/db"
	"job_board/models"
)

var database *gorm.DB
var adminName string
var adminEmail string
var adminPassword string
var adminPicture string
var adminMobileNumber string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	adminName = os.Getenv("ADMIN_NAME")
	adminEmail = os.Getenv("ADMIN_EMAIL")
	adminPassword = os.Getenv("ADMIN_PASSWORD")
	adminPicture = os.Getenv("ADMIN_PICTURE")
	adminMobileNumber = os.Getenv("ADMIN_MOBILE_NUMBER")

	if adminName == "" || adminEmail == "" || adminPassword == "" || adminPicture == "" || adminMobileNumber == "" {
		log.Fatal("Error loading super admin details")
	}
	database = db.GetDB()
	createAdmin()
}

func createAdmin() {
	log.Print("checking admin")

	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existingUser := &models.User{}
	result := tx.Where(&models.User{RoleName: models.SuperAdminRole}).First(existingUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			user := models.User{
				Name:         adminName,
				Email:        adminEmail,
				Picture:      adminPicture,
				Password:     adminPassword,
				MobileNumber: &adminMobileNumber,
				RoleName:     models.SuperAdminRole, // Ensure the role is set to SuperAdminRole
			}

			if err := tx.Create(&user).Error; err != nil {
				tx.Rollback()
				log.Panicf("failed to create admin user: %s", err)
			}

			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
				log.Printf("error committing transaction: %v", err)
				return
			}
			log.Print("created admin")
			return
		}
		tx.Rollback()
		log.Printf("error fetching admin user: %v", result.Error)
		return
	}
}

func GetSingleUser(filter models.User) (*models.User, error) {
	var user models.User
	if err := database.Preload("Profile").Preload("Companies").Preload("JobApplications").Where(&filter).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsers(filter FilterDetails, pageNumber string, pageSize string) ([]models.User, int64, int, int, error) {
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

	if err := db.Preload("Profile").Preload("Companies").Preload("JobApplications").Limit(perPage).Offset(offset).Find(&users).Error; err != nil {
		log.Println("Error finding books:", err)
		return nil, 0, 0, 0, err
	}

	return users, total, page, perPage, nil
}

func DeleteSingleUser(userID uuid.UUID) error {
	result := database.Delete(&models.User{}, userID)
	if result.RowsAffected == 0 {
		return errors.New("user already deleted")
	}
	return result.Error
}

func UpdateSingleUser(id uuid.UUID, values interface{}) (*models.User, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existingUser := &models.User{}
	result := tx.Preload("Profile").Preload("Companies").Preload("JobApplications").Where("id = ?", id).First(existingUser)
	if result.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error fetching user: %w", result.Error)
	}

	err := tx.Model(existingUser).Updates(values).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return existingUser, nil
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
	result := tx.Preload("Profile").Preload("Companies").Preload("JobApplications").Where("mobile_number = ? OR email = ?", user.MobileNumber, user.Email).First(existingUser)
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
