package job

type Request struct {
	Name string `json:"name" binding:"required"`
}