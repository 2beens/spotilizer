package services

// TODO: maybe move this to some constants or config file
const (
	SpotifyApiURL           = "https://api.spotify.com"
	URLCurrentUserPlaylists = "/v1/me/playlists"
)

var Users = NewUserService()
var UserPlaylist UserPlaylistService = SpotifyUserPlaylistService{spotifyApiURL: SpotifyApiURL, urlCurrentUserPlaylists: URLCurrentUserPlaylists}
