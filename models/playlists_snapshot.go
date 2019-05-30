package models

import "time"

// PlaylistsSnapshot is an object representing the playlist snapshot in time
type PlaylistsSnapshot struct {
	Username  string       `json:"username"`
	Timestamp time.Time    `json:"timestamp"`
	Playlists []SpPlaylist `json:"playlists"`
}
