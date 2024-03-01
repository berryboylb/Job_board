package notifications

type Subscriber struct {
	SubscriberID string                 `json:"subscriberId"`
	Name         string                 `json:"name"`
	Email        string                 `json:"email"`
	Avatar       string                 `json:"avatar"`
	Data         map[string]interface{} `json:"data"`
}

type Trigger struct {
	SubscriberID string                 `json:"subscriberId"`
	EventID      string                 `json:"eventId"`
	Name         string                 `json:"name"`
	Email        string                 `json:"email"`
	Title        string                 `json:"title"`
	Logo         string                 `json:"logo"`
	Data         map[string]interface{} `json:"data"`
}

type TriggerTopic struct {
	TopicKey string `json:"topic_key"`
	EventID  string `json:"eventId"`
	Title    string `json:"title"`
	Logo     string `json:"logo"`
}
