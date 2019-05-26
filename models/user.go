package models

type User struct {
	Username  string
	ID        string
	Auth      SpotifyAuthOptions
	FavTracks []SpAddedTrack
}
