package model

type LikeStatus struct {
	Likes     int64 `json:"likes"`
	UserLikes int64 `json:"userLikes"`
}
