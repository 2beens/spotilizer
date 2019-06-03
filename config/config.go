package config

var spotifyAPIURL = "https://api.spotify.com"
var urlCurrentUserPlaylists = "/v1/me/playlists"
var urlCurrentUserSavedTracks = "/v1/me/tracks"
var urlCurrentUser = "/v1/me"

type Config struct {
	SpotifyAPIURL             string
	URLCurrentUserPlaylists   string
	URLCurrentUserSavedTracks string
	URLCurrentUser            string
}

var Conf = &Config{
	SpotifyAPIURL:             spotifyAPIURL,
	URLCurrentUserPlaylists:   urlCurrentUserPlaylists,
	URLCurrentUserSavedTracks: urlCurrentUserSavedTracks,
	URLCurrentUser:            urlCurrentUser,
}
