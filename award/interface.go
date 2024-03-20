package award
import (
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Year        int    `json:"year" binding:"required,validYear"` //custom validation fucntion
}


func validYear(fl validator.FieldLevel) bool {
    year := fl.Field().Int()

    // Check if the year falls within a reasonable range
    // For example, consider years between 1900 and 2100 as valid
    return year >= 1900 && year <= 2100
}

type Search struct {
	Title       string `json:"title" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	Year        int    `json:"year" binding:"omitempty"`
}
