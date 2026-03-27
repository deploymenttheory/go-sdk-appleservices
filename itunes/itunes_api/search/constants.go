package search

// Media type values for SearchOptions.Media.
const (
	MediaAll        = "all"
	MediaMusic      = "music"
	MediaPodcast    = "podcast"
	MediaMusicVideo = "musicVideo"
	MediaAudiobook  = "audiobook"
	MediaShortFilm  = "shortFilm"
	MediaTVShow     = "tvShow"
	MediaMovie      = "movie"
	MediaSoftware   = "software"
	MediaEBook      = "ebook"
)

// Entity values for music searches (SearchOptions.Media = "music").
const (
	EntityMusicArtist    = "musicArtist"
	EntityMusicTrack     = "musicTrack"
	EntityAlbum          = "album"
	EntityMusicVideo     = "musicVideo"
	EntityMix            = "mix"
	EntitySong           = "song"
	EntityAllArtist      = "allArtist"
	EntityAllTrack       = "allTrack"
)

// Entity values for movie searches (SearchOptions.Media = "movie").
const (
	EntityMovieArtist = "movieArtist"
	EntityMovie       = "movie"
)

// Entity values for podcast searches (SearchOptions.Media = "podcast").
const (
	EntityPodcastAuthor  = "podcastAuthor"
	EntityPodcast        = "podcast"
)

// Entity values for software searches (SearchOptions.Media = "software").
const (
	EntitySoftware      = "software"
	EntityIPadSoftware  = "iPadSoftware"
	EntityMacSoftware   = "macSoftware"
)

// Entity values for TV show searches (SearchOptions.Media = "tvShow").
const (
	EntityTVEpisode = "tvEpisode"
	EntityTVSeason  = "tvSeason"
)

// Entity values for audiobook searches (SearchOptions.Media = "audiobook").
const (
	EntityAudiobook       = "audiobook"
	EntityAudiobookAuthor = "audiobookAuthor"
)

// Explicit content filter values for SearchOptions.Explicit.
const (
	ExplicitYes = "Yes"
	ExplicitNo  = "No"
)

// Sort values for LookupOptions.Sort.
const (
	SortRecent = "recent"
	SortPopular = "popular"
)

// Attribute values for music attribute searches (SearchOptions.Attribute).
const (
	AttributeMixTerm         = "mixTerm"
	AttributeGenreTerm       = "genreTerm"
	AttributeArtistTerm      = "artistTerm"
	AttributeComposerTerm    = "composerTerm"
	AttributeAlbumTerm       = "albumTerm"
	AttributeRatingIndex     = "ratingIndex"
	AttributeSongTerm        = "songTerm"
	AttributeTitleTerm       = "titleTerm"
	AttributeKeywordsTerm    = "keywordsTerm"
	AttributeDescriptionTerm = "descriptionTerm"
	AttributeAuthorTerm      = "authorTerm"
	AttributeDirectorTerm    = "directorTerm"
	AttributeProducerTerm    = "producerTerm"
	AttributeActorTerm       = "actorTerm"
)
