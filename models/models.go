package models

type SpGetCurrentPlaylistsResp struct {
	Href     string       `json:"href"`
	Items    []SpPlaylist `json:"items"`
	Limit    int          `json:"limit"`
	Next     string       `json:"next"`
	Offset   int          `json:"offset"`
	Previous string       `json:"previous"`
	Total    int          `json:"total"`
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
