package services

import (
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
