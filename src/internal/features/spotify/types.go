package spotify

type Platform string

const (
	PlatformSpotify  Platform = "spotify"
	PlatformJellyfin Platform = "jellyfin"
)

// Track matches your TS 'Track' type
type Track struct {
	TrackURL      *string `json:"trackUrl"` // pointer allows null
	TrackName     string  `json:"trackName"`
	AlbumName     string  `json:"albumName"`
	ArtistName    string  `json:"artistName"`
	ProgressMs    int     `json:"progressMs"`
	DurationMs    int     `json:"durationMs"`
	AlbumImageURL *string `json:"albumImageUrl"` // pointer allows null
}

// CurrentlyPlaying matches your TS 'CurrentlyPlaying' type
type CurrentlyPlaying struct {
	Item      *Track   `json:"item"` // nil if nothing playing
	Platform  Platform `json:"platform"`
	IsPlaying bool     `json:"isPlaying"`
}
