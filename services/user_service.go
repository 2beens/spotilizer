package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	c "github.com/2beens/spotilizer/config"
	"github.com/2beens/spotilizer/constants"
	"github.com/2beens/spotilizer/db"
	m "github.com/2beens/spotilizer/models"
)

type UserService struct {
	cookiesDB        db.CookiesDBClient
	usersDB          db.UsersDBClient
	username2userMap map[string]*m.User
	// TODO: add cookie expiration mechanism
	cookieID2usernameMap map[string]string
}

func NewUserService(cookiesDBClient db.CookiesDBClient, usersDB db.UsersDBClient) *UserService {
	var us UserService
	us.cookiesDB = cookiesDBClient
	us.usersDB = usersDB

	us.SyncWithDB()
	us.cookieID2usernameMap = make(map[string]string)
	us.cookieID2usernameMap = us.cookiesDB.GetCookiesInfo()
	return &us
}

func (us *UserService) AddUserCookie(cookieID string, username string) {
	us.cookieID2usernameMap[cookieID] = username
}

func (us *UserService) RemoveUserCookie(cookieID string) {
	delete(us.cookieID2usernameMap, cookieID)
	log.Println(" > user cookie removed: " + cookieID)
}

func (us *UserService) GetCookieIDByUsername(username string) (string, error) {
	for c, un := range us.cookieID2usernameMap {
		if un == username && len(c) > 0 {
			return c, nil
		}
		if len(c) == 0 {
			log.Printf(" >>> warning: found an empty cookie for user: [%s]\n", un)
		}
	}
	return "", errors.New("cookie ID not found by username")
}

func (us *UserService) GetUsernameByCookieID(cookieID string) (username string, found bool) {
	username, found = us.cookieID2usernameMap[cookieID]
	return
}

func (us *UserService) GetUserByCookieID(cookieID string) (user *m.User, err error) {
	username, found := us.cookieID2usernameMap[cookieID]
	if !found || !us.Exists(username) {
		log.Printf(" >>> error, cannot find user by cookie ID: %s\n", cookieID)
		return nil, errors.New("cannot find user by provided cookie ID")
	}
	user, _ = us.Get(username)
	return user, nil
}

func (us *UserService) SyncWithDB() {
	us.username2userMap = make(map[string]*m.User)
	// get all users from Redis
	for _, u := range *us.usersDB.GetAllUsers() {
		us.username2userMap[u.Username] = &u
		log.Printf(" > found and added user: %s\n", u.Username)
	}
}

func (us *UserService) StoreCookiesToDB() {
	us.cookiesDB.SaveCookiesInfo(us.cookieID2usernameMap)
}

func (us *UserService) Exists(username string) (found bool) {
	_, found = us.username2userMap[username]
	return
}

func (us *UserService) Get(username string) (user *m.User, err error) {
	if !us.Exists(username) {
		return nil, errors.New("cannot find user with provided ID")
	}
	user = us.username2userMap[username]
	err = nil
	return
}

func (us *UserService) Add(user *m.User) {
	us.username2userMap[user.Username] = user
	us.usersDB.SaveUser(user)
}

func (us *UserService) Save(user *m.User) (stored bool) {
	return us.usersDB.SaveUser(user)
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

func (us *UserService) GetUserByRequestCookieID(r *http.Request) (user *m.User, err error) {
	cookieID, err := r.Cookie(constants.CookieUserIDKey)
	if err != nil {
		log.Printf(" >>> %s\n", fmt.Sprintf(" >>> error, cannot find user by cookieID: %s", err.Error()))
		return
	}

	user, err = us.GetUserByCookieID(cookieID.Value)
	if err != nil {
		log.Printf(" >>> %s\n", fmt.Sprintf(" >>> cannot find user by cookie. error: %s", err.Error()))
		return
	}
	return
}
