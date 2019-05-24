package services

import (
	"errors"
	m "github.com/2beens/spotilizer/models"
)

type UserService struct {
	User2authOptionsMap map[string]m.SpotifyAuthOptions
}

func NewUserService() *UserService {
	var us UserService
	us.User2authOptionsMap = make(map[string]m.SpotifyAuthOptions)
	return &us
}

func (us *UserService) UserExists(userID string) (found bool) {
	_, found = us.User2authOptionsMap[userID]
	return
}

func (us *UserService) Get(userID string) (ao m.SpotifyAuthOptions, err error) {
	if !us.UserExists(userID) {
		return m.SpotifyAuthOptions{}, errors.New("cannot find user auth options")
	}
	ao = us.User2authOptionsMap[userID]
	err = nil
	return
}
