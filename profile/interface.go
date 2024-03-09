package profile

import (
	"github.com/google/uuid"
)

type ProfileDto struct {
	Bio                      string    `json:"bio" binding:"required"`
	Resume                   string    `json:"resume_link" binding:"required"`
	GenderID                 uuid.UUID `json:"gender_id" binding:"required"`
	CurrentSalary            float64   `json:"current_salary" binding:"omitempty"`
	CurrentSalaryCurrencyID  uuid.UUID `json:"current_salary_id" binding:"omitempty"`
	ExpectedSalary           float64   `json:"expected_salary" binding:"omitempty"`
	ExpectedSalaryCurrencyID uuid.UUID `json:"expected_salary_id" binding:"omitempty"`
}


type FilterProfileDto struct {

}