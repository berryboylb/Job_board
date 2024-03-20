package education

import (
	// "errors"
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

func createEducation(education models.Education) (*models.Education, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&education).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating a new education: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &education, nil
}

func getEducation(filter SearchEduction, pageSize string, pageNumber string) ([]models.Education, int64, int, int, error) {
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
	db := database.Model(&models.Education{})

	// Add WHERE clauses only if the corresponding filter fields are not empty or zero
	if filter.InstitutionName != "" {
		db = db.Where("institution_name LIKE ?", "%"+filter.InstitutionName+"%")
	}
	if filter.FieldOFStudy != "" {
		db = db.Where("field_of_study LIKE ?", "%"+filter.FieldOFStudy+"%")
	}
	if filter.DegreeID != uuid.Nil {
		db = db.Where("degree_id = ?", filter.DegreeID)
	}
	if filter.AcademicRankingID != uuid.Nil {
		db = db.Where("academic_ranking_id = ?", filter.AcademicRankingID)
	}
	if filter.GraduationYear != 0 {
		db = db.Where("graduation_year = ?", filter.GraduationYear)
	}
	if !filter.StartDate.IsZero() {
		// Convert time to string in the format expected by your database
		startDateStr := filter.StartDate.Format("2006-01-02")
		db = db.Where("start_date >= ?", startDateStr)
	}

	if !filter.EndDate.IsZero() {
		// Convert time to string in the format expected by your database
		endDateStr := filter.EndDate.Format("2006-01-02")
		db = db.Where("end_date <= ?", endDateStr)
	}

	// Count total number of profiles
	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Println("Error counting profiles:", err)
		return nil, 0, 0, 0, err
	}

	// Retrieve profiles with preloaded associations
	var data []models.Education
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&data).Error; err != nil {
		log.Println("Error finding education:", err)
		return nil, 0, 0, 0, err
	}

	return data, total, page, perPage, nil
}

func getSingleEducation(search models.Education) (*models.Education, error) {
	var education models.Education
	if err := database.
		Where(&search).
		First(&education).Error; err != nil {
		return nil, err
	}
	return &education, nil
}