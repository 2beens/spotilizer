package models

type DTOPlaylistSnapshot struct {
	Timestamp int64         `json:"timestamp"`
	Playlists []DTOPlaylist `json:"playlists"`
}

type DTOPlaylist struct {
	URI        string     `json:"uri"`
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	TracksHref string     `json:"trakcsHref"`
	Tracks     []DTOTrack `json:"tracks"`
}

type DTOTrack struct {
	Artists     []DTOArtist `json:"artists"`
	URI         string      `json:"uri"`
	ID          string      `json:"id"`
	TrackNumber int         `json:"track_number"`
	DurationMs  int         `json:"duration_ms"`
	Name        string      `json:"name"`
}

type DTOArtist struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
