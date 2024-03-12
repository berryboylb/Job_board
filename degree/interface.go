package degree

type Degree struct {
	Name string `json:"name" binding:"required"`
}
