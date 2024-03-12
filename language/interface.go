package language

type Language struct {
	Name string `json:"name" binding:"required"`
}