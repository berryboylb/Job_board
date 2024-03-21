package job

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"job_board/db"
	"job_board/models"
)

var database *gorm.DB

func init() {
	database = db.GetDB()
}

func createLevelFunc(Level models.Level) (*models.Level, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&Level).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("Level with the same name already exists")
		}
		return nil, fmt.Errorf("error creating a new Level: %v", err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &Level, nil
}

func getLevelFunc(name string, pageSize string, pageNumber string) ([]models.Level, int64, int, int, error) {
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
	db := database.Model(&models.Level{})

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
	var profiles []models.Level
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		log.Println("Error finding Level:", err)
		return nil, 0, 0, 0, err
	}

	return profiles, total, page, perPage, nil
}

func getSingleLevelFunc(ID uuid.UUID) (*models.Level, error) {
	Level := models.Level{}
	if err := database.
		First(&Level, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &Level, nil
}

func updateLevelFunc(ID uuid.UUID, name string) (*models.Level, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var existingRecord models.Level
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

func deleteSingle(LevelID uuid.UUID) error {
	result := database.Delete(&models.Level{}, LevelID)
	if result.RowsAffected == 0 {
		return errors.New("Level already deleted")
	}
	return result.Error
}

func createJobType(JobType models.JobType) (*models.JobType, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&JobType).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("JobType with the same name already exists")
		}
		return nil, fmt.Errorf("error creating a new JobType: %v", err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &JobType, nil
}

func getJobType(name string, pageSize string, pageNumber string) ([]models.JobType, int64, int, int, error) {
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
	db := database.Model(&models.JobType{})

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
	var profiles []models.JobType
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		log.Println("Error finding JobType:", err)
		return nil, 0, 0, 0, err
	}

	return profiles, total, page, perPage, nil
}

func getSingleJobType(ID uuid.UUID) (*models.JobType, error) {
	JobType := models.JobType{}
	if err := database.
		First(&JobType, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &JobType, nil
}

func updateJobType(ID uuid.UUID, name string) (*models.JobType, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var existingRecord models.JobType
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

func deleteSingleJobType(JobTypeID uuid.UUID) error {
	result := database.Delete(&models.JobType{}, JobTypeID)
	if result.RowsAffected == 0 {
		return errors.New("JobType already deleted")
	}
	return result.Error
}

/*job services start here*/

func createJob(Job models.Job) (*models.Job, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&Job).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating a new Job: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &Job, nil
}

func getJob(filter JobRequest, pageSize string, pageNumber string) ([]models.Job, int64, int, int, error) {
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
	db := database.Model(&models.Job{})

	if filter.Title != "" {
		db = db.Where("title LIKE ?", "%"+filter.Title+"%")
	}
	if filter.Description != "" {
		db = db.Where("description = ?", filter.Description)
	}

	if filter.CountryID != uuid.Nil {
		db = db.Where("country_id = ?", filter.CountryID)
	}

	if filter.JobTypeID != uuid.Nil {
		db = db.Where("job_type_id = ?", filter.JobTypeID)
	}

	if filter.LevelID != uuid.Nil {
		db = db.Where("level_id = ?", filter.LevelID)
	}
	if filter.CompanyID != uuid.Nil {
		db = db.Where("company_id = ?", filter.CompanyID)
	}

	if len(filter.Skills) > 0 {
		db = db.Where("skills IN (?)", pq.StringArray(filter.Skills))
	}
	if filter.Salary != 0.0 {
		db = db.Where("salary <= ?", filter.Salary)
	}

	// Count total number of profiles
	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Println("Error counting jobs:", err)
		return nil, 0, 0, 0, err
	}

	// Retrieve profiles with preloaded associations
	var data []models.Job
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&data).Error; err != nil {
		log.Println("Error finding Job:", err)
		return nil, 0, 0, 0, err
	}

	return data, total, page, perPage, nil
}

func getSingleJob(ID uuid.UUID, user models.User) (*models.Job, error) {
	var record models.Job
	if err := database.
		First(&record, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func updateJob(ID uuid.UUID, user models.User, updates JobRequest) (*models.Job, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingRecord models.Job
	if err := tx.First(&existingRecord, "id = ?", ID).Error; err != nil {
		return nil, err // Record not found or other database error
	}

	// Check if the user has permission to update the record
	if user.RoleName == models.PosterRole && existingRecord.UserID != user.ID {
		return nil, fmt.Errorf("you don't have permission to update this record")
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

func deleteSingleJob(ID uuid.UUID, user models.User) error {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingRecord models.Job
	if err := tx.First(&existingRecord, "id = ?", ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if user.RoleName == models.PosterRole && existingRecord.UserID != user.ID {
		return fmt.Errorf("you don't have permission to update this record")
	}

	result := tx.Delete(&models.Job{}, "id = ?", ID)
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

func createJobApplication(JobApplication models.JobApplication, user models.User) (*models.JobApplication, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&JobApplication).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating a new application: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &JobApplication, nil
}

func getApplications(jobID uuid.UUID, user models.User, pageSize string, pageNumber string) ([]models.JobApplication, int64, int, int, error) {
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

	var record models.Job
	if err := database.
		First(&record, "id = ?", jobID).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	if record.UserID != user.ID {
		return nil, 0, 0, 0, errors.New("you did not create this job")
	}

	db := database.Model(&models.JobApplication{})
	db = db.Where("job_id = ?", jobID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Println("Error counting profiles:", err)
		return nil, 0, 0, 0, err
	}

	offset := (page - 1) * perPage
	// Retrieve profiles with preloaded associations
	var data []models.JobApplication
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&data).Error; err != nil {
		log.Println("Error finding JobApplication:", err)
		return nil, 0, 0, 0, err
	}

	return data, total, page, perPage, nil
}

func getJobApplication(filter SearchApplication, pageSize string, pageNumber string) ([]models.JobApplication, int64, int, int, error) {
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
	db := database.Model(&models.JobApplication{})

	if filter.JobID  != uuid.Nil {
		db = db.Where("job_id = ?", filter.JobID)
	}
	if filter.ApplicantID  != uuid.Nil {
		db = db.Where("applicant_id = ?", filter.ApplicantID)
	}
	if filter.Status != "" {
		db = db.Where("status LIKE ?", "%"+filter.Status+"%")
	}
	if !filter.AppliedAt.IsZero() {
		db = db.Where("applied_at >= ?", filter.AppliedAt)
	}

	// Count total number of profiles
	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Println("Error counting profiles:", err)
		return nil, 0, 0, 0, err
	}

	// Retrieve profiles with preloaded associations
	var data []models.JobApplication
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&data).Error; err != nil {
		log.Println("Error finding JobApplication:", err)
		return nil, 0, 0, 0, err
	}

	return data, total, page, perPage, nil
}

func getSingleJobApplication(ID uuid.UUID, user models.User) (*models.JobApplication, error) {
	var record models.JobApplication
	if err := database.
		Preload("Job").
		Preload("Applicant").
		First(&record, "id = ?", ID).Error; err != nil {
		return nil, err
	}

	if user.RoleName != models.AdminRole && user.RoleName != models.SuperAdminRole {
		if record.Job.UserID != user.ID && record.ApplicantID != user.ID {
			return nil, fmt.Errorf("you don't have permission to update this record")
		}
	}

	return &record, nil
}

func updateJobApplication(ID uuid.UUID, user models.User, status models.Status) (*models.JobApplication, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingRecord models.JobApplication
	if err := tx.Preload("Job").First(&existingRecord, "id = ?", ID).Error; err != nil {
		return nil, err // Record not found or other database error
	}

	// Check if the user has permission to update the record
	if existingRecord.Job.UserID != user.ID {
		return nil, fmt.Errorf("you don't have permission to update this record")
	}

	// Update the record with the provided updates
	if err := tx.Model(&existingRecord).Update("status", status).Error; err != nil {
		return nil, err // Error updating the record
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &existingRecord, nil
}

func deleteSingleJobApplication(ID uuid.UUID, user models.User) error {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingRecord models.JobApplication
	if err := tx.First(&existingRecord, "id = ?", ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if user.RoleName != models.AdminRole || user.RoleName != models.SuperAdminRole {
		return fmt.Errorf("you don't have permission to update this record")
	}

	result := tx.Delete(&models.JobApplication{}, "id = ?", ID)
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
