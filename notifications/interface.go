package notifications

type Subscriber struct {
	Name   string                 `json:"name"`
	Email  string                 `json:"email"`
	Avatar string                 `json:"avatar"`
	Data   map[string]interface{} `json:"data"`
}

type Trigger struct {
	SubscriberID string `json:"subscriberId"`
	EventID      string `json:"eventId"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Title        string `json:"title"`
	Logo         string `json:"logo"`
}
