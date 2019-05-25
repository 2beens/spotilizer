package config

var config *Config = nil
var spotifyApiURL = "https://api.spotify.com"
var urlCurrentUserPlaylists = "/v1/me/playlists"
var urlCurrentUserSavedTracks = "/v1/me/tracks"
var urlCurrentUser = "/v1/me"

type Config struct {
	SpotifyApiURL             string
	URLCurrentUserPlaylists   string
	URLCurrentUserSavedTracks string
	URLCurrentUser            string
}

func Get() *Config {
	if config == nil {
		config = &Config{
			SpotifyApiURL:             spotifyApiURL,
			URLCurrentUserPlaylists:   urlCurrentUserPlaylists,
			URLCurrentUserSavedTracks: urlCurrentUserSavedTracks,
			URLCurrentUser:            urlCurrentUser,
		}
	}
	return config
}
