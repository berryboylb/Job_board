package ranking


type Ranking struct {
	Name string `json:"name" binding:"required"`
}