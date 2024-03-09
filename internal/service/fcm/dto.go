package fcm

type ConfigFCM struct {
	Link   string
	WebApi string
}

type CloudMessage struct {
	To               string       `json:"to"`
	Notification     Notification `json:"notification"`
	Data             any          `json:"data"`
	ContentAvailable bool         `json:"content_available"`
}

type Notification struct {
	Title            string `json:"title"`
	Body             string `json:"body"`
	Image            string `json:"image"`
	MutableContent   bool   `json:"mutable_content"`
	Sound            string `json:"sound"`
	ContentAvailable bool   `json:"content_available"`
}

type FCMResponse struct {
	MulticastID  int64       `json:"multicast_id"`
	Success      int         `json:"success"`
	Failure      int         `json:"failure"`
	CanonicalIDs int         `json:"canonical_ids"`
	Results      []FCMResult `json:"results"`
}

type FCMResult struct {
	MessageID string  `json:"message_id"`
	Error     *string `json:"error"`
}

// @super-admin

type SuperAdminSendNotification struct {
	Title  *string `json:"title" form:"title"`
	Body   *string `json:"body" form:"body"`
	UserID *int64  `json:"user_id" form:"user_id"`
}
