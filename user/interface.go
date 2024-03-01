package user

type UserDetails struct {
	Name         string `json:"name" schema:"name" `
	Email        string `json:"email" schema:"email" `
	Picture      string `json:"picture" schema:"picture"`
	RoleName     string `json:"role_name" schema:"role_name"`
	ProviderID   string `json:"provider_id" schema:"provider_id"`
	SubscriberID string `json:"subscriber_id" schema:"subscriber_id"`
	MobileNumber string `json:"mobile_number" schema:"mobile_number"`
}
