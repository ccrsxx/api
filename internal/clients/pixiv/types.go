package pixiv

import (
	"encoding/json"
	"fmt"
	"time"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    Body   `json:"body"`
}

type Body struct {
	Works []json.RawMessage `json:"works"`
	Total int               `json:"total"`
}

type FlexibleID string

func (f *FlexibleID) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexibleID(s)
		return nil
	}

	var n json.Number

	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("invalid id: %s", data)
	}

	*f = FlexibleID(n.String())

	return nil
}

type Artwork struct {
	ID             FlexibleID    `json:"id"`
	URL            string        `json:"url"`
	Title          string        `json:"title"`
	UserID         FlexibleID    `json:"userId"`
	UserName       string        `json:"userName"`
	UserAvatar     string        `json:"profileImageUrl"`
	Pages          int           `json:"pageCount"`
	XRestrict      XRestrict     `json:"xRestrict"`
	SanityLevel    SanityLevel   `json:"sl"`
	CreateDate     time.Time     `json:"createDate"` // used for user atom feeds
	AIType         AIType        `json:"aiType"`
	BookmarkData   *BookmarkData `json:"bookmarkData"`
	IsBookmarkable bool          `json:"isBookmarkable"`
	IllustType     int           `json:"illustType"`
	Tags           []string      `json:"tags"`     // used by core/popular_search
	SeriesID       string        `json:"seriesId"` // used by core/mangaseries
	SeriesTitle    string        `json:"seriesTitle"`
	Width          int           `json:"width"`
	Height         int           `json:"height"`
}

type BookmarkVisibility string

const (
	BookmarkVisibilityPublic  BookmarkVisibility = "show"
	BookmarkVisibilityPrivate BookmarkVisibility = "hide"
)

// SanityLevel represents pixiv's content rating system for artworks.
// It is more reliable and granular for authorization control than XRestrict.
//
// SanityLevel values:
//
//	0: Unreviewed - Typically seen on newly uploaded works
//	2: Safe       - Reviewed and unrestricted content
//	4: R-15       - Reviewed, mild age restriction
//	6: R-18/R-18G - Reviewed, strict age restriction
//	                (Maps to XRestrict values 1 and 2 respectively)
//
// Notes:
//   - Content with SanityLevel > 4 requires user authorization, but
//     appear to be intermittently enforced by the API.
//   - Novel routes lack SanityLevel data.
type SanityLevel int

const (
	SLUnreviewed SanityLevel = 0
	SLSafe       SanityLevel = 2
	SLR15        SanityLevel = 4
	SLR18        SanityLevel = 6
)

// pixiv returns 0, 1, 2 to filter SFW and/or NSFW artworks.
// Those values are saved in `XRestrict`.
//
// Note the hyphen in the canonical string representation;
// Go does not allow hyphens in identifiers.
type XRestrict int

const (
	Safe XRestrict = 0
	R18  XRestrict = 1
	R18G XRestrict = 2
	All  XRestrict = 3 // All is a custom value to represent all ratings.
)

// pixiv returns 0, 1, 2 to filter SFW and/or NSFW artworks..
// Those values are saved in `aiType`.
type AIType int

const (
	Unrated        AIType = 0
	NotAIGenerated AIType = 1
	AIGenerated    AIType = 2
)

// BookmarkData is a custom type to handle the following API response formats:
//
// Type 1, bookmarked:
//
//	"bookmarkData": {
//	  "id": "1234",
//	  "private": false
//	},
//
// Type 2, not bookmarked:
//
//	"bookmarkData": null
type BookmarkData struct {
	ID      string `json:"id"`
	Private bool   `json:"private"`
}

// pixiv returns 0, 1, 2 to indicate the type of illustration.
// Those values are saved in `illustType`.
type IllustType int

const (
	Illustration IllustType = 0
	Manga        IllustType = 1
	Ugoira       IllustType = 2
	Novels       IllustType = 3 // Novels is a custom value to represent novels.

)
