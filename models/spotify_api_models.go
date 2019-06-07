package models

import (
	"time"
)

type SpGetCurrentPlaylistsResp struct {
	Href     string       `json:"href"`
	Items    []SpPlaylist `json:"items"`
	Limit    int          `json:"limit"`
	Next     string       `json:"next"`
	Offset   int          `json:"offset"`
	Previous string       `json:"previous"`
	Total    int          `json:"total"`
}

// TODO: these "SpGet**" struct responses are basically the same, the only diff being the items
// see if those can be merged into one type (maybe somehow by using []interface{} for items)

type SpGetPlaylistTracksResp struct {
	Href     string            `json:"href"`
	Items    []SpPlaylistTrack `json:"items"`
	Limit    int               `json:"limit"`
	Next     string            `json:"next"`
	Offset   int               `json:"offset"`
	Previous string            `json:"previous"`
	Total    int               `json:"total"`
}

type SpGetSavedTracksResp struct {
	Href     string         `json:"href"`
	Items    []SpAddedTrack `json:"items"`
	Limit    int            `json:"limit"`
	Next     string         `json:"next"`
	Offset   int            `json:"offset"`
	Previous string         `json:"previous"`
	Total    int            `json:"total"`
}

type SpAddedTrack struct {
	AddedAt time.Time `json:"added_at"`
	Track   SpTrack   `json:"track"`
}

type SpTrack struct {
	Album            SpAlbum       `json:"album"`
	Artists          []SpArtist    `json:"artists"`
	AvailableMarkets []string      `json:"available_markets"`
	DiscNumber       int           `json:"disc_number"`
	DurationMs       int           `json:"duration_ms"`
	Explicit         bool          `json:"explicit"`
	ExternalIds      SpExternalIds `json:"external_ids"`
	ExternalUrls     SpURL         `json:"external_urls"`
	Href             string        `json:"href"`
	ID               string        `json:"id"`
	IsLocal          bool          `json:"is_local"`
	Name             string        `json:"name"`
	Popularity       int           `json:"popularity"`
	PreviewURL       string        `json:"preview_url"`
	TrackNumber      int           `json:"track_number"`
	Type             string        `json:"type"`
	URI              string        `json:"uri"`
}

type SpAlbum struct {
	AlbumType            string     `json:"album_type"`
	Artists              []SpArtist `json:"artists"`
	AvailableMarkets     []string   `json:"available_markets"`
	ExternalUrls         SpURL      `json:"external_urls"`
	Href                 string     `json:"href"`
	ID                   string     `json:"id"`
	Images               []SpImage  `json:"images"`
	Name                 string     `json:"name"`
	ReleaseDate          string     `json:"release_date"`
	ReleaseDatePrecision string     `json:"release_date_precision"`
	TotalTracks          int        `json:"total_tracks"`
	Type                 string     `json:"type"`
	URI                  string     `json:"uri"`
}

type SpPlaylist struct {
	Collaborative bool        `json:"collaborative"`
	ExternalUrls  SpURL       `json:"external_urls"`
	Href          string      `json:"href"`
	ID            string      `json:"id"`
	Images        []SpImage   `json:"images"`
	Name          string      `json:"name"`
	Owner         SpUser      `json:"owner"`
	PrimaryColor  interface{} `json:"primary_color"`
	Public        bool        `json:"public"`
	SnapshotID    string      `json:"snapshot_id"`
	Tracks        SpTracks    `json:"tracks"`
	Type          string      `json:"type"`
	URI           string      `json:"uri"`
}

type SpExternalIds struct {
	Isrc string `json:"isrc"`
}

type SpTracks struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}

type SpURL struct {
	Spotify string `json:"spotify"`
}

type SpImage struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type SpArtist struct {
	ExternalUrls SpURL  `json:"external_urls"`
	Href         string `json:"href"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}

type SpError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SpAPIError struct {
	Error SpError `json:"error"`
}

type SpUser struct {
	Birthdate       string                   `json:"birthdate"`
	Country         string                   `json:"country"`
	DisplayName     string                   `json:"display_name"`
	Email           string                   `json:"email"`
	ExplicitContent SpExplicitContentOptions `json:"explicit_content"`
	ExternalUrls    SpURL                    `json:"external_urls"`
	Followers       SpUserFollowers          `json:"followers"`
	Href            string                   `json:"href"`
	ID              string                   `json:"id"`
	Images          []SpImage                `json:"images"`
	Product         string                   `json:"product"`
	Type            string                   `json:"type"`
	URI             string                   `json:"uri"`
}

type SpExplicitContentOptions struct {
	FilterEnabled bool `json:"filter_enabled"`
	FilterLocked  bool `json:"filter_locked"`
}

type SpUserFollowers struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}

type SpPlaylistTrack struct {
	AddedAt        time.Time        `json:"added_at"`
	AddedBy        SpAddedBy        `json:"added_by"`
	IsLocal        bool             `json:"is_local"`
	PrimaryColor   interface{}      `json:"primary_color"`
	Track          SpTrack          `json:"track"`
	VideoThumbnail SpVideoThumbnail `json:"video_thumbnail"`
}

type SpAddedBy struct {
	ExternalUrls SpExternalUrls `json:"external_urls"`
	Href         string         `json:"href"`
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	URI          string         `json:"uri"`
}

type SpExternalUrls struct {
	Spotify string `json:"spotify"`
}

type SpVideoThumbnail struct {
	URL interface{} `json:"url"`
}
