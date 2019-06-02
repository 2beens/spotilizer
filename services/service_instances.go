package services

// TODO: this just somehow does not seem the best way to do it - keeping an instances of services here
// 		 gotta think about this a bit later

import "github.com/2beens/spotilizer/db"

var Users *UserService
var UserPlaylist *SpotifyUserPlaylistService

func InitServices() {
	Users = NewUserService(db.GetCookiesDBClient(), db.GetUsersDBClient())
	UserPlaylist = NewSpotifyUserPlaylistService()
}
