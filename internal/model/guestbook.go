package model

import "time"

type Guestbook struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}
