package model

type Platform string

const (
	PlatformSpotify  Platform = "spotify"
	PlatformJellyfin Platform = "jellyfin"
)

type Track struct {
	TrackURL      *string `json:"trackUrl"`
	TrackName     string  `json:"trackName"`
	AlbumName     string `json:"albumName"`
	ArtistName    string  `json:"artistName"`
	ProgressMs    int     `json:"progressMs"`
	DurationMs    int     `json:"durationMs"`
	AlbumImageURL *string `json:"albumImageUrl"`
}

type CurrentlyPlaying struct {
	Item      *Track   `json:"item"`
	Platform  Platform `json:"platform"`
	IsPlaying bool     `json:"isPlaying"`
}

func NewDefaultCurrentlyPlaying(platform Platform) CurrentlyPlaying {
	return CurrentlyPlaying{
		Platform:  platform,
		IsPlaying: false,
		Item:      nil,
	}
}
