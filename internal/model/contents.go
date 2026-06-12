package model

import "time"

type Content struct {
	Slug      string     `json:"slug"`
	Type      string     `json:"type"`
	Views     int64      `json:"views,omitempty"`
	Likes     int64      `json:"likes,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
