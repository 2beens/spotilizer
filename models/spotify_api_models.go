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

// TODO: these two responses are basically the same, the only diff being the items
// see if those can be merged into one type (maybe somehow by using []interface{} for items)

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
	ExternalUrls     SpUrl         `json:"external_urls"`
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
	ExternalUrls         SpUrl      `json:"external_urls"`
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
	ExternalUrls  SpUrl       `json:"external_urls"`
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

type SpUrl struct {
	Spotify string `json:"spotify"`
}

type SpImage struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type SpUser struct {
	DisplayName  string `json:"display_name"`
	ExternalUrls SpUrl  `json:"external_urls"`
	Href         string `json:"href"`
	ID           string `json:"id"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}

type SpArtist struct {
	ExternalUrls SpUrl  `json:"external_urls"`
	Href         string `json:"href"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	URI          string `json:"uri"`
}

type SpError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}