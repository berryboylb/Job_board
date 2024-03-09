package gender

type Gender struct {
	Name string `json:"name" binding:"required"`
}



type searchGender struct {
	Name string `json:"name" binding:"omitempty"`
}