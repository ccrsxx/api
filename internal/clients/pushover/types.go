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
	Title     string          `json:"title,omitempty"`
	Device    string          `json:"device,omitempty"`
	URL       string          `json:"url,omitempty"`
	URLTitle  string          `json:"url_title,omitempty"`
	Priority  MessagePriority `json:"priority,omitempty"`
	Sound     string          `json:"sound,omitempty"`
	Timestamp int64           `json:"timestamp,omitempty"`
	TTL       int             `json:"ttl,omitempty"`

	// HTML and Monospace are mutually exclusive; set one or neither (not both).
	HTML      int `json:"html,omitempty"`
	Monospace int `json:"monospace,omitempty"`

	// Required only when Priority == MessagePriorityEmergency
	Retry    int    `json:"retry,omitempty"`
	Expire   int    `json:"expire,omitempty"`
	Callback string `json:"callback,omitempty"`
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
