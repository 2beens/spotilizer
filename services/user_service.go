package services

import (
	"encoding/json"
	"errors"
	c "github.com/2beens/spotilizer/config"
	db "github.com/2beens/spotilizer/db"
	m "github.com/2beens/spotilizer/models"
	"log"
)

type UserService struct {
	id2userMap map[string]*m.User
}

func NewUserService() *UserService {
	var us UserService
	us.SyncWithDB()
	return &us
}

func (us *UserService) SyncWithDB() {
	us.id2userMap = make(map[string]*m.User)
	// get all users from Redis
	for _, u := range *db.GetAllUsers() {
		us.id2userMap[u.ID] = &u
		log.Printf(" > found and added user: %s\n", u.Username)
	}
}

func (us *UserService) Exists(userID string) (found bool) {
	_, found = us.id2userMap[userID]
	return
}

func (us *UserService) Get(userID string) (user *m.User, err error) {
	if !us.Exists(userID) {
		return nil, errors.New("cannot find user with provided ID")
	}
	user = us.id2userMap[userID]
	err = nil
	return
}

func (us *UserService) Add(user *m.User) {
	us.id2userMap[user.ID] = user
	db.SaveUser(user)
}

func (us *UserService) GetByUsername(username string) (u *m.User) {
	for _, u := range us.id2userMap {
		if u.Username == username {
			return u
		}
	}
	return nil //, errors.New("cannot find user with provided username")
}

func (us *UserService) GetUserFromSpotify(ao *m.SpotifyAuthOptions) (user *m.SpUser, err error) {
	body, err := getFromSpotify(c.Get().SpotifyApiURL, c.Get().URLCurrentUser, ao)
	if err != nil {
		log.Printf(" >>> error getting current user playlists. details: %v\n", err)
		return nil, err
	}
	json.Unmarshal(body, &user)
	return user, nil
}
