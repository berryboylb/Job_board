package language

import (
	"github.com/google/uuid"
)
type Language struct {
	Name string `json:"name" binding:"required"`
}

type Request struct {
	Name string `json:"name" binding:"required"`
	LanguageID  uuid.UUID  `json:"language_id" binding:"required"`
	LanguageProficiencyID uuid.UUID `json:"language_proficiency_id" binding:"required"`
}

type Search struct {
	Name string `json:"name" binding:"omitempty"`
	LanguageID  uuid.UUID  `json:"language_id" binding:"omitempty"`
	LanguageProficiencyID uuid.UUID `json:"language_proficiency_id" binding:"omitempty"`
}