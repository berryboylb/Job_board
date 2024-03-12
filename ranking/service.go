package ranking


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

func createAcademicRanking(academicRanking models.AcademicRanking) (*models.AcademicRanking, error) {
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&academicRanking).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("AcademicRanking with the same name already exists")
		}
		return nil, fmt.Errorf("error creating a new AcademicRanking: %v", err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &academicRanking, nil
}

func getAcademicRanking(name string, pageSize string, pageNumber string) ([]models.AcademicRanking, int64, int, int, error) {
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
	db := database.Model(&models.AcademicRanking{})

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
	var profiles []models.AcademicRanking
	if err := db.
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		log.Println("Error finding AcademicRanking:", err)
		return nil, 0, 0, 0, err
	}

	return profiles, total, page, perPage, nil
}

func getSingleAcademicRanking(search models.AcademicRanking) (*models.AcademicRanking, error) {
	AcademicRanking := models.AcademicRanking{}
	if err := database.
		Where(&search).
		First(&AcademicRanking).Error; err != nil {
		return nil, err
	}
	return &AcademicRanking, nil
}

func updateAcademicRanking(academicRankingID uuid.UUID, name string) (*models.AcademicRanking, error) {
    tx := database.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    result := tx.Model(&models.AcademicRanking{}).Where("id = ?", academicRankingID).Update("name", name)
    if result.Error != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error updating AcademicRanking: %w", result.Error)
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Return the updated AcademicRanking
    updatedAcademicRanking := &models.AcademicRanking{ID: academicRankingID, Name: name}
    return updatedAcademicRanking, nil
}

func deleteSingle(academicRankingID uuid.UUID) error {
	result := database.Delete(&models.AcademicRanking{}, academicRankingID)
	if result.RowsAffected == 0 {
		return errors.New("AcademicRanking already deleted")
	}
	return result.Error
}