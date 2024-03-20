package project

import (
	"fmt"
	"time"
)

type Request struct {
	ProjectName string `json:"project_name" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date" binding:"omitempty"`
}

var format = "2006-01-02"

func (r *Request) ParseStartDate() (time.Time, error) {
	return time.Parse(format, r.StartDate)
}

// ParseEndDate parses the EndDate string into a time.Time object
func (r *Request) ParseEndDate() (time.Time, error) {
	if r.EndDate == nil {
		return time.Time{}, nil
	}
	return time.Parse(format, *r.EndDate)
}

func (r *Request) ValidateDates() (*time.Time, *time.Time, error) {
	startDate, err := r.ParseStartDate()
	if err != nil {
		return nil, nil, err
	}

	endDate, err := r.ParseEndDate()
	if err != nil {
		return nil, nil, err
	}

	// Check if both StartDate and EndDate are provided before comparing
	if !startDate.IsZero() && !endDate.IsZero() && !startDate.Before(endDate) {
		return nil, nil, fmt.Errorf("StartDate must be before EndDate")
	}

	return &startDate, &endDate, nil
}

type Search struct {
	ProjectName string    `json:"project_name" binding:"omitempty"`
	Title       string    `json:"title" binding:"omitempty"`
	Description string    `json:"description" binding:"omitempty"`
	StartDate   time.Time `json:"start_date" binding:"omitempty"`
	EndDate     time.Time `json:"end_date" binding:"omitempty"`
}
