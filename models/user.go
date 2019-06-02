package models

import (
	"fmt"
)

// User is an object representing the user of this service, not Spotify
type User struct {
	Username  string
	Auth      *SpotifyAuthOptions
	FavTracks *[]SpAddedTrack
	Playlists *[]SpPlaylist
}

func (u User) String() string {
	tracksLen := "<nil>"
	playlistsLen := "<nil>"
	if u.FavTracks != nil {
		tracksLen = fmt.Sprintf("%d", len(*u.FavTracks))
	}
	if u.Playlists != nil {
		playlistsLen = fmt.Sprintf("%d", len(*u.Playlists))
	}
	return fmt.Sprintf("[%s]: tracks [%s], playlists [%s], auth: [%v]",
		u.Username, tracksLen, playlistsLen, *u.Auth)
}
