package model

type Statistic struct {
	Type       string `json:"type"`
	TotalPosts int64  `json:"totalPosts"`
	TotalViews int64  `json:"totalViews"`
	TotalLikes int64  `json:"totalLikes"`
}
