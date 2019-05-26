package models

// User is an object representing the user of this service, not Spotify
type User struct {
	Username  string
	ID        string
	Auth      *SpotifyAuthOptions
	FavTracks *[]SpAddedTrack
}
