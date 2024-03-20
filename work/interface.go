package work


type Request struct {
	Attachment string `json:"attachment" binding:"required"`
	Link       string `json:"link" binding:"required"`
	Description string `json:"description" binding:"required"`
}



type Search struct {
	Attachment string    `json:"attachment" binding:"omitempty"`
	Link       string    `json:"link" binding:"omitempty"`
	Description string    `json:"description" binding:"omitempty"`
}