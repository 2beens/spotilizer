package services

// TODO: this just somehow does not seem the best way to do it - keeping an instances of services here
// 		 gotta think about this a bit later

var Users = NewUserService()
var UserPlaylist UserPlaylistService = NewSpotifyUserPlaylistService()
