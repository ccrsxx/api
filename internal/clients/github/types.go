package github

import "time"

type User struct {
	Login             string    `json:"login"`
	ID                int64     `json:"id"`
	UserViewType      *string   `json:"user_view_type,omitempty"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        *string   `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              *string   `json:"name"`
	Company           *string   `json:"company"`
	Blog              *string   `json:"blog"`
	Location          *string   `json:"location"`
	Email             *string   `json:"email"`
	NotificationEmail *string   `json:"notification_email,omitempty"`
	Hireable          *bool     `json:"hireable"`
	Bio               *string   `json:"bio"`
	TwitterUsername   *string   `json:"twitter_username,omitempty"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Plan              *Plan     `json:"plan,omitempty"`
	PrivateGists      *int      `json:"private_gists,omitempty"`
	TotalPrivateRepos *int      `json:"total_private_repos,omitempty"`
	OwnedPrivateRepos *int      `json:"owned_private_repos,omitempty"`
	DiskUsage         *int      `json:"disk_usage,omitempty"`
	Collaborators     *int      `json:"collaborators,omitempty"`
}

type Plan struct {
	Collaborators int    `json:"collaborators"`
	Name          string `json:"name"`
	Space         int    `json:"space"`
	PrivateRepos  int    `json:"private_repos"`
}
