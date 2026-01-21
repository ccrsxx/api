package spotify

// CurrentlyPlayingType maps to CurrentlyPlayingResponse.CurrentlyPlayingType
type CurrentlyPlayingType string

const (
	CurrentlyPlayingTypeTrack   CurrentlyPlayingType = "track"
	CurrentlyPlayingTypeEpisode CurrentlyPlayingType = "episode"
	CurrentlyPlayingTypeAd      CurrentlyPlayingType = "ad"
	CurrentlyPlayingTypeUnknown CurrentlyPlayingType = "unknown"
)

// SpotifyContextType maps to ContextObject.Type
type SpotifyContextType string

const (
	ContextTypeArtist   SpotifyContextType = "artist"
	ContextTypePlaylist SpotifyContextType = "playlist"
	ContextTypeAlbum    SpotifyContextType = "album"
	ContextTypeShow     SpotifyContextType = "show"
)

// SpotifyCurrentlyPlaying maps to SpotifyApi.CurrentlyPlayingResponse
type SpotifyCurrentlyPlaying struct {
	Timestamp            int64                `json:"timestamp"`
	ProgressMs           int                  `json:"progress_ms"`
	IsPlaying            bool                 `json:"is_playing"`
	CurrentlyPlayingType CurrentlyPlayingType `json:"currently_playing_type"`
	Item                 *SpotifyItem         `json:"item"`
	Context              *SpotifyContext      `json:"context"`
	Device               *SpotifyDevice       `json:"device"`
	Actions              *SpotifyActions      `json:"actions,omitempty"` // Added commonly used "Disallows"
}

// SpotifyItem is a SUPERSET struct handling both TrackObjectFull and EpisodeObject
type SpotifyItem struct {
	// Shared Fields
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Type         CurrentlyPlayingType `json:"type"` // "track" or "episode"
	DurationMs   int                  `json:"duration_ms"`
	Explicit     bool                 `json:"explicit"`
	URI          string               `json:"uri"`
	ExternalURLs SpotifyExternalURLs  `json:"external_urls"`
	Href         string               `json:"href"`

	// --------------------------------------------------------
	// Track Specific (Present if Type == "track")
	// --------------------------------------------------------
	Album      *SpotifyAlbum   `json:"album,omitempty"`
	Artists    []SpotifyArtist `json:"artists,omitempty"`
	PreviewURL string          `json:"preview_url,omitempty"`
	IsLocal    bool            `json:"is_local,omitempty"`
	Popularity int             `json:"popularity,omitempty"` // Added useful missing field

	// --------------------------------------------------------
	// Episode Specific (Present if Type == "episode")
	// --------------------------------------------------------
	Show            *SpotifyShow        `json:"show,omitempty"`
	Images          []SpotifyImage      `json:"images,omitempty"` // Episodes have images directly
	Description     string              `json:"description,omitempty"`
	HTMLDescription string              `json:"html_description,omitempty"`
	ReleaseDate     string              `json:"release_date,omitempty"` // Episodes have release dates
	ResumePoint     *SpotifyResumePoint `json:"resume_point,omitempty"` // Crucial for Podcasts
}

// SpotifyResumePoint (New Struct for Episodes)
type SpotifyResumePoint struct {
	FullyPlayed      bool `json:"fully_played"`
	ResumePositionMs int  `json:"resume_position_ms"`
}

// SpotifyActions (New Struct for "actions" / "disallows")
type SpotifyActions struct {
	Disallows map[string]bool `json:"disallows"`
}

// SpotifyAlbum maps to AlbumObjectSimplified
type SpotifyAlbum struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	AlbumType    string              `json:"album_type"`
	Images       []SpotifyImage      `json:"images"`
	ReleaseDate  string              `json:"release_date"`
	ExternalURLs SpotifyExternalURLs `json:"external_urls"`
}

// SpotifyArtist maps to ArtistObjectSimplified
type SpotifyArtist struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Type         string              `json:"type"`
	ExternalURLs SpotifyExternalURLs `json:"external_urls"`
}

// SpotifyShow maps to ShowObjectSimplified
type SpotifyShow struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Publisher    string              `json:"publisher"`
	Images       []SpotifyImage      `json:"images"`
	ExternalURLs SpotifyExternalURLs `json:"external_urls"`
}

// SpotifyContext maps to ContextObject
type SpotifyContext struct {
	Type         SpotifyContextType  `json:"type"`
	HREF         string              `json:"href"`
	ExternalURLs SpotifyExternalURLs `json:"external_urls"`
	URI          string              `json:"uri"`
}

// SpotifyDevice maps to UserDevice
type SpotifyDevice struct {
	ID            string `json:"id"`
	IsActive      bool   `json:"is_active"`
	IsPrivate     bool   `json:"is_private_session"`
	IsRestricted  bool   `json:"is_restricted"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	VolumePercent int    `json:"volume_percent"`
}

// Helpers
type SpotifyImage struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

// SpotifyExternalURLs maps to ExternalURL object
type SpotifyExternalURLs struct {
	Spotify string `json:"spotify"`
}
