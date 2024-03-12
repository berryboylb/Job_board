package company

type Request struct {
	Name string `json:"name" binding:"required"`
}