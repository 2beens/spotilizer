package models

import "time"

// PlaylistsSnapshot is an object representing the playlist snapshot in time
type PlaylistsSnapshot struct {
	Username  string             `json:"username"`
	Timestamp time.Time          `json:"timestamp"`
	Playlists []PlaylistSnapshot `json:"playlists"`
}

type PlaylistSnapshot struct {
	Playlist SpPlaylist
	Tracks   []SpPlaylistTrack
}

// FavTracksSnapshot is an object representing the list of favorite saved tracks of a user
type FavTracksSnapshot struct {
	Username  string         `json:"username"`
	Timestamp time.Time      `json:"timestamp"`
	Tracks    []SpAddedTrack `json:"tracks"`
}
