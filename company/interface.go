package company

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Request struct {
	Name string `json:"name" binding:"required"`
}

type CompanyRequest struct {
	Name            string    `json:"name" binding:"required"`
	Website         string    `json:"website" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	IndustryID      uuid.UUID `json:"start_date" binding:"required"`
	Established     string    `json:"end_date" binding:"required"`
	Location        string    `json:"location" binding:"required"`
	EmployeesSizeID uuid.UUID `json:"employees_size_id" binding:"required"`
	Logo            string    `json:"logo" binding:"required"`
}

var format = "2006-01-02"

func (c *CompanyRequest) ParseStartDate() (time.Time, error) {
	// Parse the Established date string into a time.Time object
	establishedTime, err := time.Parse(format, c.Established)
	if err != nil {
		return time.Time{}, err // Return error if parsing fails
	}

	// Validate the Established date if needed
	// For example, ensure that the Established date is not in the future
	if establishedTime.After(time.Now()) {
		return time.Time{}, fmt.Errorf("Established date cannot be in the future")
	}

	return establishedTime, nil
}

type SearchCompanyRequest struct {
	Name            string    `json:"name" binding:"required"`
	Website         string    `json:"website" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	IndustryID      uuid.UUID `json:"start_date" binding:"required"`
	Established     time.Time    `json:"end_date" binding:"required"`
	Location        string    `json:"location" binding:"required"`
	EmployeesSizeID uuid.UUID `json:"employees_size_id" binding:"required"`
	Logo            string    `json:"logo" binding:"required"`
}
