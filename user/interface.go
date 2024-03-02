package user

type UserDetails struct {
	Name         string `json:"name" binding:"omitempty"`
	Email        string `json:"email" binding:"omitempty,email"`
	Picture      string `json:"picture" binding:"omitempty"`
	// RoleName     string `json:"role_name" binding:"omitempty"`
	// ProviderID   string `json:"provider_id" binding:"omitempty"`
	// SubscriberID string `json:"subscriber_id" binding:"omitempty"`
	MobileNumber string `json:"mobile_number" binding:"omitempty"`
}

type FilterDetails struct {
	Name         string `json:"name" binding:"omitempty"`
	Email        string `json:"email" binding:"omitempty,email"`
	Picture      string `json:"picture" binding:"omitempty"`
	RoleName     string `json:"role_name" binding:"omitempty"`
	ProviderID   string `json:"provider_id" binding:"omitempty"`
	SubscriberID string `json:"subscriber_id" binding:"omitempty"`
	MobileNumber string `json:"mobile_number" binding:"omitempty"`
}

type Admin struct {
    Name         string `json:"name" binding:"required"`
    Email        string `json:"email" binding:"required,email"`
    Picture      string `json:"picture" binding:"omitempty"`
    MobileNumber string `json:"mobile_number" binding:"required"`
}

