package jellyfin

// MediaType Enum
type MediaType string

const (
	MediaTypeUnknown MediaType = "Unknown"
	MediaTypeVideo   MediaType = "Video"
	MediaTypeAudio   MediaType = "Audio"
	MediaTypePhoto   MediaType = "Photo"
	MediaTypeBook    MediaType = "Book"
)

// GeneralCommandType Enum
type GeneralCommandType string

const (
	CommandMoveUp                 GeneralCommandType = "MoveUp"
	CommandMoveDown               GeneralCommandType = "MoveDown"
	CommandMoveLeft               GeneralCommandType = "MoveLeft"
	CommandMoveRight              GeneralCommandType = "MoveRight"
	CommandPageUp                 GeneralCommandType = "PageUp"
	CommandPageDown               GeneralCommandType = "PageDown"
	CommandPreviousLetter         GeneralCommandType = "PreviousLetter"
	CommandNextLetter             GeneralCommandType = "NextLetter"
	CommandToggleOsd              GeneralCommandType = "ToggleOsd"
	CommandToggleContextMenu      GeneralCommandType = "ToggleContextMenu"
	CommandSelect                 GeneralCommandType = "Select"
	CommandBack                   GeneralCommandType = "Back"
	CommandTakeScreenshot         GeneralCommandType = "TakeScreenshot"
	CommandSendKey                GeneralCommandType = "SendKey"
	CommandSendString             GeneralCommandType = "SendString"
	CommandGoHome                 GeneralCommandType = "GoHome"
	CommandGoToSettings           GeneralCommandType = "GoToSettings"
	CommandVolumeUp               GeneralCommandType = "VolumeUp"
	CommandVolumeDown             GeneralCommandType = "VolumeDown"
	CommandMute                   GeneralCommandType = "Mute"
	CommandUnmute                 GeneralCommandType = "Unmute"
	CommandToggleMute             GeneralCommandType = "ToggleMute"
	CommandSetVolume              GeneralCommandType = "SetVolume"
	CommandSetAudioStreamIndex    GeneralCommandType = "SetAudioStreamIndex"
	CommandSetSubtitleStreamIndex GeneralCommandType = "SetSubtitleStreamIndex"
	CommandToggleFullscreen       GeneralCommandType = "ToggleFullscreen"
	CommandDisplayContent         GeneralCommandType = "DisplayContent"
	CommandGoToSearch             GeneralCommandType = "GoToSearch"
	CommandDisplayMessage         GeneralCommandType = "DisplayMessage"
	CommandSetRepeatMode          GeneralCommandType = "SetRepeatMode"
	CommandChannelUp              GeneralCommandType = "ChannelUp"
	CommandChannelDown            GeneralCommandType = "ChannelDown"
	CommandGuide                  GeneralCommandType = "Guide"
	CommandToggleStats            GeneralCommandType = "ToggleStats"
	CommandPlayMediaSource        GeneralCommandType = "PlayMediaSource"
	CommandPlayTrailers           GeneralCommandType = "PlayTrailers"
	CommandSetShuffleQueue        GeneralCommandType = "SetShuffleQueue"
	CommandPlayState              GeneralCommandType = "PlayState"
	CommandPlayNext               GeneralCommandType = "PlayNext"
	CommandToggleOsdMenu          GeneralCommandType = "ToggleOsdMenu"
	CommandPlay                   GeneralCommandType = "Play"
	CommandSetMaxStreamingBitrate GeneralCommandType = "SetMaxStreamingBitrate"
	CommandSetPlaybackOrder       GeneralCommandType = "SetPlaybackOrder"
)

// PlayMethod Enum
type PlayMethod string

const (
	PlayMethodTranscode    PlayMethod = "Transcode"
	PlayMethodDirectStream PlayMethod = "DirectStream"
	PlayMethodDirectPlay   PlayMethod = "DirectPlay"
)

// RepeatMode Enum
type RepeatMode string

const (
	RepeatNone RepeatMode = "RepeatNone"
	RepeatAll  RepeatMode = "RepeatAll"
	RepeatOne  RepeatMode = "RepeatOne"
)

// PlaybackOrder Enum
type PlaybackOrder string

const (
	PlaybackOrderDefault PlaybackOrder = "Default"
	PlaybackOrderShuffle PlaybackOrder = "Shuffle"
)

// HardwareAccelerationType Enum
type HardwareAccelerationType string

const (
	HardwareAccelNone         HardwareAccelerationType = "none"
	HardwareAccelAmf          HardwareAccelerationType = "amf"
	HardwareAccelQsv          HardwareAccelerationType = "qsv"
	HardwareAccelNvenc        HardwareAccelerationType = "nvenc"
	HardwareAccelV4l2m2m      HardwareAccelerationType = "v4l2m2m"
	HardwareAccelVaapi        HardwareAccelerationType = "vaapi"
	HardwareAccelVideotoolbox HardwareAccelerationType = "videotoolbox"
	HardwareAccelRkmpp        HardwareAccelerationType = "rkmpp"
)

// TranscodeReason Enum
type TranscodeReason string

const (
	ReasonContainerNotSupported        TranscodeReason = "ContainerNotSupported"
	ReasonVideoCodecNotSupported       TranscodeReason = "VideoCodecNotSupported"
	ReasonAudioCodecNotSupported       TranscodeReason = "AudioCodecNotSupported"
	ReasonSubtitleCodecNotSupported    TranscodeReason = "SubtitleCodecNotSupported"
	ReasonAudioIsExternal              TranscodeReason = "AudioIsExternal"
	ReasonSecondaryAudioNotSupported   TranscodeReason = "SecondaryAudioNotSupported"
	ReasonVideoProfileNotSupported     TranscodeReason = "VideoProfileNotSupported"
	ReasonVideoLevelNotSupported       TranscodeReason = "VideoLevelNotSupported"
	ReasonVideoResolutionNotSupported  TranscodeReason = "VideoResolutionNotSupported"
	ReasonVideoBitDepthNotSupported    TranscodeReason = "VideoBitDepthNotSupported"
	ReasonVideoFramerateNotSupported   TranscodeReason = "VideoFramerateNotSupported"
	ReasonRefFramesNotSupported        TranscodeReason = "RefFramesNotSupported"
	ReasonAnamorphicVideoNotSupported  TranscodeReason = "AnamorphicVideoNotSupported"
	ReasonInterlacedVideoNotSupported  TranscodeReason = "InterlacedVideoNotSupported"
	ReasonAudioChannelsNotSupported    TranscodeReason = "AudioChannelsNotSupported"
	ReasonAudioProfileNotSupported     TranscodeReason = "AudioProfileNotSupported"
	ReasonAudioSampleRateNotSupported  TranscodeReason = "AudioSampleRateNotSupported"
	ReasonAudioBitDepthNotSupported    TranscodeReason = "AudioBitDepthNotSupported"
	ReasonContainerBitrateExceedsLimit TranscodeReason = "ContainerBitrateExceedsLimit"
	ReasonVideoBitrateNotSupported     TranscodeReason = "VideoBitrateNotSupported"
	ReasonAudioBitrateNotSupported     TranscodeReason = "AudioBitrateNotSupported"
	ReasonUnknownVideoStreamInfo       TranscodeReason = "UnknownVideoStreamInfo"
	ReasonUnknownAudioStreamInfo       TranscodeReason = "UnknownAudioStreamInfo"
	ReasonDirectPlayError              TranscodeReason = "DirectPlayError"
	ReasonVideoRangeTypeNotSupported   TranscodeReason = "VideoRangeTypeNotSupported"
	ReasonVideoCodecTagNotSupported    TranscodeReason = "VideoCodecTagNotSupported"
	ReasonStreamCountExceedsLimit      TranscodeReason = "StreamCountExceedsLimit"
)

// BaseItemKind Enum
type BaseItemKind string

const (
	KindAggregateFolder       BaseItemKind = "AggregateFolder"
	KindAudio                 BaseItemKind = "Audio"
	KindAudioBook             BaseItemKind = "AudioBook"
	KindBasePluginFolder      BaseItemKind = "BasePluginFolder"
	KindBook                  BaseItemKind = "Book"
	KindBoxSet                BaseItemKind = "BoxSet"
	KindChannel               BaseItemKind = "Channel"
	KindChannelFolderItem     BaseItemKind = "ChannelFolderItem"
	KindCollectionFolder      BaseItemKind = "CollectionFolder"
	KindEpisode               BaseItemKind = "Episode"
	KindFolder                BaseItemKind = "Folder"
	KindGenre                 BaseItemKind = "Genre"
	KindManualPlaylistsFolder BaseItemKind = "ManualPlaylistsFolder"
	KindMovie                 BaseItemKind = "Movie"
	KindLiveTvChannel         BaseItemKind = "LiveTvChannel"
	KindLiveTvProgram         BaseItemKind = "LiveTvProgram"
	KindMusicAlbum            BaseItemKind = "MusicAlbum"
	KindMusicArtist           BaseItemKind = "MusicArtist"
	KindMusicGenre            BaseItemKind = "MusicGenre"
	KindMusicVideo            BaseItemKind = "MusicVideo"
	KindPerson                BaseItemKind = "Person"
	KindPhoto                 BaseItemKind = "Photo"
	KindPhotoAlbum            BaseItemKind = "PhotoAlbum"
	KindPlaylist              BaseItemKind = "Playlist"
	KindPlaylistsFolder       BaseItemKind = "PlaylistsFolder"
	KindProgram               BaseItemKind = "Program"
	KindRecording             BaseItemKind = "Recording"
	KindSeason                BaseItemKind = "Season"
	KindSeries                BaseItemKind = "Series"
	KindStudio                BaseItemKind = "Studio"
	KindTrailer               BaseItemKind = "Trailer"
	KindTvChannel             BaseItemKind = "TvChannel"
	KindTvProgram             BaseItemKind = "TvProgram"
	KindUserRootFolder        BaseItemKind = "UserRootFolder"
	KindUserView              BaseItemKind = "UserView"
	KindVideo                 BaseItemKind = "Video"
	KindYear                  BaseItemKind = "Year"
)

// PersonKind Enum
type PersonKind string

const (
	PersonUnknown     PersonKind = "Unknown"
	PersonActor       PersonKind = "Actor"
	PersonDirector    PersonKind = "Director"
	PersonComposer    PersonKind = "Composer"
	PersonWriter      PersonKind = "Writer"
	PersonGuestStar   PersonKind = "GuestStar"
	PersonProducer    PersonKind = "Producer"
	PersonConductor   PersonKind = "Conductor"
	PersonLyricist    PersonKind = "Lyricist"
	PersonArranger    PersonKind = "Arranger"
	PersonEngineer    PersonKind = "Engineer"
	PersonMixer       PersonKind = "Mixer"
	PersonRemixer     PersonKind = "Remixer"
	PersonCreator     PersonKind = "Creator"
	PersonArtist      PersonKind = "Artist"
	PersonAlbumArtist PersonKind = "AlbumArtist"
	PersonAuthor      PersonKind = "Author"
	PersonIllustrator PersonKind = "Illustrator"
	PersonPenciller   PersonKind = "Penciller"
	PersonInker       PersonKind = "Inker"
	PersonColorist    PersonKind = "Colorist"
	PersonLetterer    PersonKind = "Letterer"
	PersonCoverArtist PersonKind = "CoverArtist"
	PersonEditor      PersonKind = "Editor"
	PersonTranslator  PersonKind = "Translator"
)

// Other Enums
type LocationType string
type IsoType string
type VideoType string
type ChannelType string
type ProgramAudio string
type PlayAccess string
type DayOfWeek string
type ImageOrientation string
type DlnaProfileType string
type CodecType string
type ProfileConditionType string
type ProfileConditionValue string
type SubtitleDeliveryMethod string
type MediaStreamProtocol string
type TranscodeSeekInfo string
type EncodingContext string
type MediaProtocol string
type MediaSourceType string
type TransportStreamTimestamp string
type VideoRange string
type VideoRangeType string
type AudioSpatialFormat string
type MediaStreamType string
type Video3DFormat string
type ExtraType string

const (
	VideoRangeSDR VideoRange = "SDR"
	VideoRangeHDR VideoRange = "HDR"
)

type ProfileCondition struct {
	Condition  ProfileConditionType  `json:"Condition"`
	Property   ProfileConditionValue `json:"Property"`
	Value      *string               `json:"Value"`
	IsRequired bool                  `json:"IsRequired"`
}

type DirectPlayProfile struct {
	Container  string          `json:"Container"`
	AudioCodec *string         `json:"AudioCodec"`
	VideoCodec *string         `json:"VideoCodec"`
	Type       DlnaProfileType `json:"Type"`
}

type TranscodingProfile struct {
	Container                 string              `json:"Container"`
	Type                      DlnaProfileType     `json:"Type"`
	VideoCodec                string              `json:"VideoCodec"`
	AudioCodec                string              `json:"AudioCodec"`
	Protocol                  MediaStreamProtocol `json:"Protocol"`
	EstimateContentLength     bool                `json:"EstimateContentLength"`
	EnableMpegtsM2TsMode      bool                `json:"EnableMpegtsM2TsMode"`
	TranscodeSeekInfo         TranscodeSeekInfo   `json:"TranscodeSeekInfo"`
	CopyTimestamps            bool                `json:"CopyTimestamps"`
	Context                   EncodingContext     `json:"Context"`
	EnableSubtitlesInManifest bool                `json:"EnableSubtitlesInManifest"`
	MaxAudioChannels          *string             `json:"MaxAudioChannels"`
	MinSegments               int                 `json:"MinSegments"`
	SegmentLength             int                 `json:"SegmentLength"`
	BreakOnNonKeyFrames       bool                `json:"BreakOnNonKeyFrames"`
	Conditions                []ProfileCondition  `json:"Conditions"`
	EnableAudioVbrEncoding    bool                `json:"EnableAudioVbrEncoding"`
}

type ContainerProfile struct {
	Type         DlnaProfileType    `json:"Type"`
	Conditions   []ProfileCondition `json:"Conditions"`
	Container    *string            `json:"Container"`
	SubContainer *string            `json:"SubContainer"`
}

type CodecProfile struct {
	Type            CodecType          `json:"Type"`
	Conditions      []ProfileCondition `json:"Conditions"`
	ApplyConditions []ProfileCondition `json:"ApplyConditions"`
	Codec           *string            `json:"Codec"`
	Container       *string            `json:"Container"`
	SubContainer    *string            `json:"SubContainer"`
}

type SubtitleProfile struct {
	Format    *string                `json:"Format"`
	Method    SubtitleDeliveryMethod `json:"Method"`
	DidlMode  *string                `json:"DidlMode"`
	Language  *string                `json:"Language"`
	Container *string                `json:"Container"`
}

type DeviceProfile struct {
	Name                             *string              `json:"Name"`
	Id                               *string              `json:"Id"`
	MaxStreamingBitrate              *int                 `json:"MaxStreamingBitrate"`
	MaxStaticBitrate                 *int                 `json:"MaxStaticBitrate"`
	MusicStreamingTranscodingBitrate *int                 `json:"MusicStreamingTranscodingBitrate"`
	MaxStaticMusicBitrate            *int                 `json:"MaxStaticMusicBitrate"`
	DirectPlayProfiles               []DirectPlayProfile  `json:"DirectPlayProfiles"`
	TranscodingProfiles              []TranscodingProfile `json:"TranscodingProfiles"`
	ContainerProfiles                []ContainerProfile   `json:"ContainerProfiles"`
	CodecProfiles                    []CodecProfile       `json:"CodecProfiles"`
	SubtitleProfiles                 []SubtitleProfile    `json:"SubtitleProfiles"`
}

type ClientCapabilities struct {
	PlayableMediaTypes           []MediaType          `json:"PlayableMediaTypes"`
	SupportedCommands            []GeneralCommandType `json:"SupportedCommands"`
	SupportsMediaControl         bool                 `json:"SupportsMediaControl"`
	SupportsPersistentIdentifier bool                 `json:"SupportsPersistentIdentifier"`
	DeviceProfile                *DeviceProfile       `json:"DeviceProfile"`
	AppStoreUrl                  *string              `json:"AppStoreUrl"`
	IconUrl                      *string              `json:"IconUrl"`
}

type MediaStream struct {
	Codec                     *string                 `json:"Codec"`
	CodecTag                  *string                 `json:"CodecTag"`
	Language                  *string                 `json:"Language"`
	ColorRange                *string                 `json:"ColorRange"`
	ColorSpace                *string                 `json:"ColorSpace"`
	ColorTransfer             *string                 `json:"ColorTransfer"`
	ColorPrimaries            *string                 `json:"ColorPrimaries"`
	DvVersionMajor            *int                    `json:"DvVersionMajor"`
	DvVersionMinor            *int                    `json:"DvVersionMinor"`
	DvProfile                 *int                    `json:"DvProfile"`
	DvLevel                   *int                    `json:"DvLevel"`
	RpuPresentFlag            *int                    `json:"RpuPresentFlag"`
	ElPresentFlag             *int                    `json:"ElPresentFlag"`
	BlPresentFlag             *int                    `json:"BlPresentFlag"`
	DvBlSignalCompatibilityId *int                    `json:"DvBlSignalCompatibilityId"`
	Rotation                  *int                    `json:"Rotation"`
	Comment                   *string                 `json:"Comment"`
	TimeBase                  *string                 `json:"TimeBase"`
	CodecTimeBase             *string                 `json:"CodecTimeBase"`
	Title                     *string                 `json:"Title"`
	Hdr10PlusPresentFlag      *bool                   `json:"Hdr10PlusPresentFlag"`
	VideoRange                VideoRange              `json:"VideoRange"`
	VideoRangeType            VideoRangeType          `json:"VideoRangeType"`
	VideoDoViTitle            *string                 `json:"VideoDoViTitle"`
	AudioSpatialFormat        AudioSpatialFormat      `json:"AudioSpatialFormat"`
	LocalizedUndefined        *string                 `json:"LocalizedUndefined"`
	LocalizedDefault          *string                 `json:"LocalizedDefault"`
	LocalizedForced           *string                 `json:"LocalizedForced"`
	LocalizedExternal         *string                 `json:"LocalizedExternal"`
	LocalizedHearingImpaired  *string                 `json:"LocalizedHearingImpaired"`
	DisplayTitle              *string                 `json:"DisplayTitle"`
	NalLengthSize             *string                 `json:"NalLengthSize"`
	IsInterlaced              bool                    `json:"IsInterlaced"`
	IsAVC                     *bool                   `json:"IsAVC"`
	ChannelLayout             *string                 `json:"ChannelLayout"`
	BitRate                   *int                    `json:"BitRate"`
	BitDepth                  *int                    `json:"BitDepth"`
	RefFrames                 *int                    `json:"RefFrames"`
	PacketLength              *int                    `json:"PacketLength"`
	Channels                  *int                    `json:"Channels"`
	SampleRate                *int                    `json:"SampleRate"`
	IsDefault                 bool                    `json:"IsDefault"`
	IsForced                  bool                    `json:"IsForced"`
	IsHearingImpaired         bool                    `json:"IsHearingImpaired"`
	Height                    *int                    `json:"Height"`
	Width                     *int                    `json:"Width"`
	AverageFrameRate          *float64                `json:"AverageFrameRate"`
	RealFrameRate             *float64                `json:"RealFrameRate"`
	ReferenceFrameRate        *float64                `json:"ReferenceFrameRate"`
	Profile                   *string                 `json:"Profile"`
	Type                      MediaStreamType         `json:"Type"`
	AspectRatio               *string                 `json:"AspectRatio"`
	Index                     int                     `json:"Index"`
	Score                     *int                    `json:"Score"`
	IsExternal                bool                    `json:"IsExternal"`
	DeliveryMethod            *SubtitleDeliveryMethod `json:"DeliveryMethod"`
	DeliveryUrl               *string                 `json:"DeliveryUrl"`
	IsExternalUrl             *bool                   `json:"IsExternalUrl"`
	IsTextSubtitleStream      bool                    `json:"IsTextSubtitleStream"`
	SupportsExternalStream    bool                    `json:"SupportsExternalStream"`
	Path                      *string                 `json:"Path"`
	PixelFormat               *string                 `json:"PixelFormat"`
	Level                     *float64                `json:"Level"`
	IsAnamorphic              *bool                   `json:"IsAnamorphic"`
}

type MediaAttachment struct {
	Codec       *string `json:"Codec"`
	CodecTag    *string `json:"CodecTag"`
	Comment     *string `json:"Comment"`
	Index       int     `json:"Index"`
	FileName    *string `json:"FileName"`
	MimeType    *string `json:"MimeType"`
	DeliveryUrl *string `json:"DeliveryUrl"`
}

type MediaSourceInfo struct {
	Protocol                            MediaProtocol             `json:"Protocol"`
	Id                                  *string                   `json:"Id"`
	Path                                *string                   `json:"Path"`
	EncoderPath                         *string                   `json:"EncoderPath"`
	EncoderProtocol                     *MediaProtocol            `json:"EncoderProtocol"`
	Type                                MediaSourceType           `json:"Type"`
	Container                           *string                   `json:"Container"`
	Size                                *int64                    `json:"Size"`
	Name                                *string                   `json:"Name"`
	IsRemote                            bool                      `json:"IsRemote"`
	ETag                                *string                   `json:"ETag"`
	RunTimeTicks                        *int64                    `json:"RunTimeTicks"`
	ReadAtNativeFramerate               bool                      `json:"ReadAtNativeFramerate"`
	IgnoreDts                           bool                      `json:"IgnoreDts"`
	IgnoreIndex                         bool                      `json:"IgnoreIndex"`
	GenPtsInput                         bool                      `json:"GenPtsInput"`
	SupportsTranscoding                 bool                      `json:"SupportsTranscoding"`
	SupportsDirectStream                bool                      `json:"SupportsDirectStream"`
	SupportsDirectPlay                  bool                      `json:"SupportsDirectPlay"`
	IsInfiniteStream                    bool                      `json:"IsInfiniteStream"`
	UseMostCompatibleTranscodingProfile bool                      `json:"UseMostCompatibleTranscodingProfile"`
	RequiresOpening                     bool                      `json:"RequiresOpening"`
	OpenToken                           *string                   `json:"OpenToken"`
	RequiresClosing                     bool                      `json:"RequiresClosing"`
	LiveStreamId                        *string                   `json:"LiveStreamId"`
	BufferMs                            *int                      `json:"BufferMs"`
	RequiresLooping                     bool                      `json:"RequiresLooping"`
	SupportsProbing                     bool                      `json:"SupportsProbing"`
	VideoType                           *VideoType                `json:"VideoType"`
	IsoType                             *IsoType                  `json:"IsoType"`
	Video3DFormat                       *Video3DFormat            `json:"Video3DFormat"`
	MediaStreams                        []MediaStream             `json:"MediaStreams"`
	MediaAttachments                    []MediaAttachment         `json:"MediaAttachments"`
	Formats                             []string                  `json:"Formats"`
	Bitrate                             *int                      `json:"Bitrate"`
	FallbackMaxStreamingBitrate         *int                      `json:"FallbackMaxStreamingBitrate"`
	Timestamp                           *TransportStreamTimestamp `json:"Timestamp"`
	RequiredHttpHeaders                 map[string]*string        `json:"RequiredHttpHeaders"`
	TranscodingUrl                      *string                   `json:"TranscodingUrl"`
	TranscodingSubProtocol              MediaStreamProtocol       `json:"TranscodingSubProtocol"`
	TranscodingContainer                *string                   `json:"TranscodingContainer"`
	AnalyzeDurationMs                   *int                      `json:"AnalyzeDurationMs"`
	DefaultAudioStreamIndex             *int                      `json:"DefaultAudioStreamIndex"`
	DefaultSubtitleStreamIndex          *int                      `json:"DefaultSubtitleStreamIndex"`
	HasSegments                         bool                      `json:"HasSegments"`
}

type ExternalUrl struct {
	Name *string `json:"Name"`
	Url  *string `json:"Url"`
}

type MediaUrl struct {
	Url  *string `json:"Url"`
	Name *string `json:"Name"`
}

type UserItemData struct {
	Rating                *float64 `json:"Rating"`
	PlayedPercentage      *float64 `json:"PlayedPercentage"`
	UnplayedItemCount     *int     `json:"UnplayedItemCount"`
	PlaybackPositionTicks *int64   `json:"PlaybackPositionTicks"`
	PlayCount             *int     `json:"PlayCount"`
	IsFavorite            bool     `json:"IsFavorite"`
	Likes                 *bool    `json:"Likes"`
	LastPlayedDate        *string  `json:"LastPlayedDate"`
	Played                bool     `json:"Played"`
	Key                   string   `json:"Key"`
	ItemId                string   `json:"ItemId"`
}

type ImageBlurHashes struct {
	Primary    map[string]string `json:"Primary"`
	Art        map[string]string `json:"Art"`
	Backdrop   map[string]string `json:"Backdrop"`
	Banner     map[string]string `json:"Banner"`
	Logo       map[string]string `json:"Logo"`
	Thumb      map[string]string `json:"Thumb"`
	Disc       map[string]string `json:"Disc"`
	Box        map[string]string `json:"Box"`
	Screenshot map[string]string `json:"Screenshot"`
	Menu       map[string]string `json:"Menu"`
	Chapter    map[string]string `json:"Chapter"`
	BoxRear    map[string]string `json:"BoxRear"`
	Profile    map[string]string `json:"Profile"`
}

type BaseItemPerson struct {
	Name            *string          `json:"Name"`
	Id              string           `json:"Id"`
	Role            *string          `json:"Role"`
	Type            PersonKind       `json:"Type"`
	PrimaryImageTag *string          `json:"PrimaryImageTag"`
	ImageBlurHashes *ImageBlurHashes `json:"ImageBlurHashes"`
}

type NameGuidPair struct {
	Name *string `json:"Name"`
	Id   string  `json:"Id"`
}

type ChapterInfo struct {
	StartPositionTicks int64   `json:"StartPositionTicks"`
	Name               *string `json:"Name"`
	ImagePath          *string `json:"ImagePath"`
	ImageDateModified  string  `json:"ImageDateModified"`
	ImageTag           *string `json:"ImageTag"`
}

type TrickplayInfo struct {
	Width          int `json:"Width"`
	Height         int `json:"Height"`
	TileWidth      int `json:"TileWidth"`
	TileHeight     int `json:"TileHeight"`
	ThumbnailCount int `json:"ThumbnailCount"`
	Interval       int `json:"Interval"`
	Bandwidth      int `json:"Bandwidth"`
}

type BaseItem struct {
	Name                         *string                             `json:"Name"`
	OriginalTitle                *string                             `json:"OriginalTitle"`
	ServerId                     *string                             `json:"ServerId"`
	Id                           string                              `json:"Id"`
	Etag                         *string                             `json:"Etag"`
	SourceType                   *string                             `json:"SourceType"`
	PlaylistItemId               *string                             `json:"PlaylistItemId"`
	DateCreated                  *string                             `json:"DateCreated"`
	DateLastMediaAdded           *string                             `json:"DateLastMediaAdded"`
	ExtraType                    *ExtraType                          `json:"ExtraType"`
	AirsBeforeSeasonNumber       *int                                `json:"AirsBeforeSeasonNumber"`
	AirsAfterSeasonNumber        *int                                `json:"AirsAfterSeasonNumber"`
	AirsBeforeEpisodeNumber      *int                                `json:"AirsBeforeEpisodeNumber"`
	CanDelete                    *bool                               `json:"CanDelete"`
	CanDownload                  *bool                               `json:"CanDownload"`
	HasLyrics                    *bool                               `json:"HasLyrics"`
	HasSubtitles                 *bool                               `json:"HasSubtitles"`
	PreferredMetadataLanguage    *string                             `json:"PreferredMetadataLanguage"`
	PreferredMetadataCountryCode *string                             `json:"PreferredMetadataCountryCode"`
	Container                    *string                             `json:"Container"`
	SortName                     *string                             `json:"SortName"`
	ForcedSortName               *string                             `json:"ForcedSortName"`
	Video3DFormat                *Video3DFormat                      `json:"Video3DFormat"`
	PremiereDate                 *string                             `json:"PremiereDate"`
	ExternalUrls                 []ExternalUrl                       `json:"ExternalUrls"`
	MediaSources                 []MediaSourceInfo                   `json:"MediaSources"`
	CriticRating                 *float64                            `json:"CriticRating"`
	ProductionLocations          []string                            `json:"ProductionLocations"`
	Path                         *string                             `json:"Path"`
	EnableMediaSourceDisplay     *bool                               `json:"EnableMediaSourceDisplay"`
	OfficialRating               *string                             `json:"OfficialRating"`
	CustomRating                 *string                             `json:"CustomRating"`
	ChannelId                    *string                             `json:"ChannelId"`
	ChannelName                  *string                             `json:"ChannelName"`
	Overview                     *string                             `json:"Overview"`
	Taglines                     []string                            `json:"Taglines"`
	Genres                       []string                            `json:"Genres"`
	CommunityRating              *float64                            `json:"CommunityRating"`
	CumulativeRunTimeTicks       *int64                              `json:"CumulativeRunTimeTicks"`
	RunTimeTicks                 *int64                              `json:"RunTimeTicks"`
	PlayAccess                   *PlayAccess                         `json:"PlayAccess"`
	AspectRatio                  *string                             `json:"AspectRatio"`
	ProductionYear               *int                                `json:"ProductionYear"`
	IsPlaceHolder                *bool                               `json:"IsPlaceHolder"`
	Number                       *string                             `json:"Number"`
	ChannelNumber                *string                             `json:"ChannelNumber"`
	IndexNumber                  *int                                `json:"IndexNumber"`
	IndexNumberEnd               *int                                `json:"IndexNumberEnd"`
	ParentIndexNumber            *int                                `json:"ParentIndexNumber"`
	RemoteTrailers               []MediaUrl                          `json:"RemoteTrailers"`
	ProviderIds                  map[string]*string                  `json:"ProviderIds"`
	IsHD                         *bool                               `json:"IsHD"`
	IsFolder                     *bool                               `json:"IsFolder"`
	ParentId                     *string                             `json:"ParentId"`
	Type                         BaseItemKind                        `json:"Type"`
	People                       []BaseItemPerson                    `json:"People"`
	Studios                      []NameGuidPair                      `json:"Studios"`
	GenreItems                   []NameGuidPair                      `json:"GenreItems"`
	ParentLogoItemId             *string                             `json:"ParentLogoItemId"`
	ParentBackdropItemId         *string                             `json:"ParentBackdropItemId"`
	ParentBackdropImageTags      []string                            `json:"ParentBackdropImageTags"`
	LocalTrailerCount            *int                                `json:"LocalTrailerCount"`
	UserData                     *UserItemData                       `json:"UserData"`
	RecursiveItemCount           *int                                `json:"RecursiveItemCount"`
	ChildCount                   *int                                `json:"ChildCount"`
	SeriesName                   *string                             `json:"SeriesName"`
	SeriesId                     *string                             `json:"SeriesId"`
	SeasonId                     *string                             `json:"SeasonId"`
	SpecialFeatureCount          *int                                `json:"SpecialFeatureCount"`
	DisplayPreferencesId         *string                             `json:"DisplayPreferencesId"`
	Status                       *string                             `json:"Status"`
	AirTime                      *string                             `json:"AirTime"`
	AirDays                      []DayOfWeek                         `json:"AirDays"`
	Tags                         []string                            `json:"Tags"`
	PrimaryImageAspectRatio      *float64                            `json:"PrimaryImageAspectRatio"`
	Artists                      []string                            `json:"Artists"`
	ArtistItems                  []NameGuidPair                      `json:"ArtistItems"`
	Album                        *string                             `json:"Album"`
	CollectionType               *string                             `json:"CollectionType"`
	DisplayOrder                 *string                             `json:"DisplayOrder"`
	AlbumId                      *string                             `json:"AlbumId"`
	AlbumPrimaryImageTag         *string                             `json:"AlbumPrimaryImageTag"`
	SeriesPrimaryImageTag        *string                             `json:"SeriesPrimaryImageTag"`
	AlbumArtist                  *string                             `json:"AlbumArtist"`
	AlbumArtists                 []NameGuidPair                      `json:"AlbumArtists"`
	SeasonName                   *string                             `json:"SeasonName"`
	MediaStreams                 []MediaStream                       `json:"MediaStreams"`
	VideoType                    *VideoType                          `json:"VideoType"`
	PartCount                    *int                                `json:"PartCount"`
	MediaSourceCount             *int                                `json:"MediaSourceCount"`
	ImageTags                    map[string]string                   `json:"ImageTags"`
	BackdropImageTags            []string                            `json:"BackdropImageTags"`
	ScreenshotImageTags          []string                            `json:"ScreenshotImageTags"`
	ParentLogoImageTag           *string                             `json:"ParentLogoImageTag"`
	ParentArtItemId              *string                             `json:"ParentArtItemId"`
	ParentArtImageTag            *string                             `json:"ParentArtImageTag"`
	SeriesThumbImageTag          *string                             `json:"SeriesThumbImageTag"`
	ImageBlurHashes              *ImageBlurHashes                    `json:"ImageBlurHashes"`
	SeriesStudio                 *string                             `json:"SeriesStudio"`
	ParentThumbItemId            *string                             `json:"ParentThumbItemId"`
	ParentThumbImageTag          *string                             `json:"ParentThumbImageTag"`
	ParentPrimaryImageItemId     *string                             `json:"ParentPrimaryImageItemId"`
	ParentPrimaryImageTag        *string                             `json:"ParentPrimaryImageTag"`
	Chapters                     []ChapterInfo                       `json:"Chapters"`
	Trickplay                    map[string]map[string]TrickplayInfo `json:"Trickplay"`
	LocationType                 *LocationType                       `json:"LocationType"`
	IsoType                      *IsoType                            `json:"IsoType"`
	MediaType                    MediaType                           `json:"MediaType"`
	EndDate                      *string                             `json:"EndDate"`
	LockedFields                 []string                            `json:"LockedFields"`
	TrailerCount                 *int                                `json:"TrailerCount"`
	MovieCount                   *int                                `json:"MovieCount"`
	SeriesCount                  *int                                `json:"SeriesCount"`
	ProgramCount                 *int                                `json:"ProgramCount"`
	EpisodeCount                 *int                                `json:"EpisodeCount"`
	SongCount                    *int                                `json:"SongCount"`
	AlbumCount                   *int                                `json:"AlbumCount"`
	ArtistCount                  *int                                `json:"ArtistCount"`
	MusicVideoCount              *int                                `json:"MusicVideoCount"`
	LockData                     *bool                               `json:"LockData"`
	Width                        *int                                `json:"Width"`
	Height                       *int                                `json:"Height"`
	CameraMake                   *string                             `json:"CameraMake"`
	CameraModel                  *string                             `json:"CameraModel"`
	Software                     *string                             `json:"Software"`
	ExposureTime                 *float64                            `json:"ExposureTime"`
	FocalLength                  *float64                            `json:"FocalLength"`
	ImageOrientation             *ImageOrientation                   `json:"ImageOrientation"`
	Aperture                     *float64                            `json:"Aperture"`
	ShutterSpeed                 *float64                            `json:"ShutterSpeed"`
	Latitude                     *float64                            `json:"Latitude"`
	Longitude                    *float64                            `json:"Longitude"`
	Altitude                     *float64                            `json:"Altitude"`
	IsoSpeedRating               *int                                `json:"IsoSpeedRating"`
	SeriesTimerId                *string                             `json:"SeriesTimerId"`
	ProgramId                    *string                             `json:"ProgramId"`
	ChannelPrimaryImageTag       *string                             `json:"ChannelPrimaryImageTag"`
	StartDate                    *string                             `json:"StartDate"`
	CompletionPercentage         *float64                            `json:"CompletionPercentage"`
	IsRepeat                     *bool                               `json:"IsRepeat"`
	EpisodeTitle                 *string                             `json:"EpisodeTitle"`
	ChannelType                  *ChannelType                        `json:"ChannelType"`
	Audio                        *ProgramAudio                       `json:"Audio"`
	IsMovie                      *bool                               `json:"IsMovie"`
	IsSports                     *bool                               `json:"IsSports"`
	IsSeries                     *bool                               `json:"IsSeries"`
	IsLive                       *bool                               `json:"IsLive"`
	IsNews                       *bool                               `json:"IsNews"`
	IsKids                       *bool                               `json:"IsKids"`
	IsPremiere                   *bool                               `json:"IsPremiere"`
	TimerId                      *string                             `json:"TimerId"`
	NormalizationGain            *float64                            `json:"NormalizationGain"`
	CurrentProgram               *BaseItem                           `json:"CurrentProgram"`
}

type PlayerStateInfo struct {
	PositionTicks       *int64         `json:"PositionTicks"`
	CanSeek             bool           `json:"CanSeek"`
	IsPaused            bool           `json:"IsPaused"`
	IsMuted             bool           `json:"IsMuted"`
	VolumeLevel         *int           `json:"VolumeLevel"`
	AudioStreamIndex    *int           `json:"AudioStreamIndex"`
	SubtitleStreamIndex *int           `json:"SubtitleStreamIndex"`
	MediaSourceId       *string        `json:"MediaSourceId"`
	PlayMethod          *PlayMethod    `json:"PlayMethod"`
	RepeatMode          *RepeatMode    `json:"RepeatMode"`
	PlaybackOrder       *PlaybackOrder `json:"PlaybackOrder"`
	LiveStreamId        *string        `json:"LiveStreamId"`
}

type SessionUserInfo struct {
	UserId   string  `json:"UserId"`
	UserName *string `json:"UserName"`
}

type TranscodingInfo struct {
	AudioCodec               *string                   `json:"AudioCodec"`
	VideoCodec               *string                   `json:"VideoCodec"`
	Container                *string                   `json:"Container"`
	IsVideoDirect            bool                      `json:"IsVideoDirect"`
	IsAudioDirect            bool                      `json:"IsAudioDirect"`
	Bitrate                  *int                      `json:"Bitrate"`
	Framerate                *float64                  `json:"Framerate"`
	CompletionPercentage     *float64                  `json:"CompletionPercentage"`
	Width                    *int                      `json:"Width"`
	Height                   *int                      `json:"Height"`
	AudioChannels            *int                      `json:"AudioChannels"`
	HardwareAccelerationType *HardwareAccelerationType `json:"HardwareAccelerationType"`
	TranscodeReasons         []TranscodeReason         `json:"TranscodeReasons"`
}

type QueueItem struct {
	Id             string  `json:"Id"`
	PlaylistItemId *string `json:"PlaylistItemId"`
}

// SessionInfo is the root struct usually returned by /Sessions
type SessionInfo struct {
	PlayState                *PlayerStateInfo     `json:"PlayState"`
	AdditionalUsers          []SessionUserInfo    `json:"AdditionalUsers"`
	Capabilities             *ClientCapabilities  `json:"Capabilities"`
	RemoteEndPoint           *string              `json:"RemoteEndPoint"`
	PlayableMediaTypes       []MediaType          `json:"PlayableMediaTypes"`
	Id                       *string              `json:"Id"`
	UserId                   string               `json:"UserId"`
	UserName                 *string              `json:"UserName"`
	Client                   *string              `json:"Client"`
	LastActivityDate         string               `json:"LastActivityDate"`
	LastPlaybackCheckIn      string               `json:"LastPlaybackCheckIn"`
	LastPausedDate           *string              `json:"LastPausedDate"`
	DeviceName               *string              `json:"DeviceName"`
	DeviceType               *string              `json:"DeviceType"`
	NowPlayingItem           *BaseItem            `json:"NowPlayingItem"`
	NowViewingItem           *BaseItem            `json:"NowViewingItem"`
	DeviceId                 *string              `json:"DeviceId"`
	ApplicationVersion       *string              `json:"ApplicationVersion"`
	TranscodingInfo          *TranscodingInfo     `json:"TranscodingInfo"`
	IsActive                 bool                 `json:"IsActive"`
	SupportsMediaControl     bool                 `json:"SupportsMediaControl"`
	SupportsRemoteControl    bool                 `json:"SupportsRemoteControl"`
	NowPlayingQueue          []QueueItem          `json:"NowPlayingQueue"`
	NowPlayingQueueFullItems []BaseItem           `json:"NowPlayingQueueFullItems"`
	HasCustomDeviceName      bool                 `json:"HasCustomDeviceName"`
	PlaylistItemId           *string              `json:"PlaylistItemId"`
	ServerId                 *string              `json:"ServerId"`
	UserPrimaryImageTag      *string              `json:"UserPrimaryImageTag"`
	SupportedCommands        []GeneralCommandType `json:"SupportedCommands"`
}
