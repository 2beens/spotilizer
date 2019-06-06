package models

import (
	"fmt"
)

// User is an object representing the user of this service, not Spotify
type User struct {
	Username string
	Auth     *SpotifyAuthOptions
}

func (u User) String() string {
	return fmt.Sprintf("[%s]: auth: [%v]", u.Username, *u.Auth)
}
