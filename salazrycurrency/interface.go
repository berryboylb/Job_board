package salazrycurrency


type SalaryCurrency struct {
	Name string `json:"name" binding:"required"`
}