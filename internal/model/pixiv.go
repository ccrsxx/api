package model

import "time"

type Bookmark struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	ArtistID    string    `json:"artistId"`
	ArtistName  string    `json:"artistName"`
	ImageURL    string    `json:"imageUrl"`
	PixivURL    string    `json:"pixivUrl"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Tags        []string  `json:"tags"`
	AiGenerated bool      `json:"aiGenerated"`
	CreatedAt   time.Time `json:"createdAt"`
}
