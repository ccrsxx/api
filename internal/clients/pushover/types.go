package pushover

// MessagePriority maps to the Pushover message priority values.
type MessagePriority int

const (
	MessagePriorityLowest    MessagePriority = -2
	MessagePriorityLow       MessagePriority = -1
	MessagePriorityNormal    MessagePriority = 0
	MessagePriorityHigh      MessagePriority = 1
	MessagePriorityEmergency MessagePriority = 2
)

// MessageRequest represents the payload sent to the Pushover messages API.
// Ref: https://pushover.net/api#messages
type MessageRequest struct {
	// Required fields
	Token   string `json:"token"`
	User    string `json:"user"`
	Message string `json:"message"`

	// Optional fields
	Title     string          `json:"title,omitzero"`
	Device    string          `json:"device,omitzero"`
	URL       string          `json:"url,omitzero"`
	URLTitle  string          `json:"url_title,omitzero"`
	Priority  MessagePriority `json:"priority,omitzero"`
	Sound     string          `json:"sound,omitzero"`
	Timestamp int64           `json:"timestamp,omitzero"`
	TTL       int             `json:"ttl,omitzero"`

	// HTML and Monospace are mutually exclusive; set one or neither (not both).
	HTML      int `json:"html,omitzero"`
	Monospace int `json:"monospace,omitzero"`

	// Required only when Priority == MessagePriorityEmergency
	Retry    int    `json:"retry,omitzero"`
	Expire   int    `json:"expire,omitzero"`
	Callback string `json:"callback,omitzero"`
}

// MessageResponse represents the response returned by the Pushover messages API.
// On success (HTTP 200): status=1, request=<uuid>.
// On error (HTTP 4xx): status=0, request=<uuid>, errors=[...].
// On emergency priority success: additionally includes receipt=<id>.
// Ref: https://pushover.net/api#response
type MessageResponse struct {
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Receipt string   `json:"receipt,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}
