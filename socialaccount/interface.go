package socialaccount

import (
	"github.com/google/uuid"
)

type Request struct {
	Link          string    `json:"link" binding:"required"`
	SocialMediaID uuid.UUID `json:"social_media_id" binding:"required"`
}

type Search struct {
	Link          string    `json:"link" binding:"omitempty"`
	SocialMediaID uuid.UUID `json:"social_media_id" binding:"omitempty"`
}

type SocialRequest struct {
	Name         string    `json:"name" binding:"required"`
}


