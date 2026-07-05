package cloudflare

import "time"

type Turnstile struct {
	Success     bool               `json:"success"`
	ErrorCodes  []string           `json:"error-codes"`
	ChallengeTs time.Time          `json:"challenge_ts"`
	Hostname    string             `json:"hostname"`
	Action      string             `json:"action"`
	CData       string             `json:"cdata"`
	Metadata    *TurnstileMetadata `json:"metadata"`
}

type TurnstileMetadata struct {
	Interactive bool   `json:"interactive"`
	EphemeralID string `json:"ephemeral_id"`
}
