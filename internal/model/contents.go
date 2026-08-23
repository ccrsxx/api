package model

import "time"

type Content struct {
	Slug      string     `json:"slug"`
	Type      string     `json:"type"`
	Views     int64     `json:"views"`
	Likes     int64     `json:"likes"`
	CreatedAt time.Time `json:"createdAt,omitzero"`
	UpdatedAt time.Time `json:"updatedAt,omitzero"`
}
