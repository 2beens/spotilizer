package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/2beens/spotilizer/config"
	"github.com/2beens/spotilizer/constants"
	"github.com/2beens/spotilizer/db"
	"github.com/2beens/spotilizer/models"
)

type UserService struct {
	cookiesDB            db.CookiesDBClient
	usersDB              db.UsersDBClient
	username2userMap     map[string]*models.User
	cookieID2usernameMap map[string]string
}

func NewUserService(cookiesDBClient db.CookiesDBClient, usersDB db.UsersDBClient) *UserService {
	us := new(UserService)
	us.cookiesDB = cookiesDBClient
	us.usersDB = usersDB

	us.SyncWithDB()
	us.cookieID2usernameMap = make(map[string]string)
	us.cookieID2usernameMap = us.cookiesDB.GetCookiesInfo()
	return us
}

func NewUserServiceTest() *UserService {
	us := new(UserService)
	us.usersDB = db.NewUsersDBTest([]models.User{})
	us.cookieID2usernameMap = make(map[string]string)
	us.username2userMap = make(map[string]*models.User)
	return us
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

func (us *UserService) GetUserByCookieID(cookieID string) (user *models.User, err error) {
	username, found := us.cookieID2usernameMap[cookieID]
	if !found || !us.Exists(username) {
		log.Printf(" >>> error, cannot find user by cookie ID: %s\n", cookieID)
		return nil, errors.New("cannot find user by provided cookie ID")
	}
	user, _ = us.Get(username)
	return user, nil
}

func (us *UserService) SyncWithDB() {
	us.username2userMap = make(map[string]*models.User)
	// get all users from Redis
	for _, u := range us.usersDB.GetAllUsers() {
		user := u
		us.username2userMap[u.Username] = &user
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

func (us *UserService) Get(username string) (user *models.User, err error) {
	if !us.Exists(username) {
		return nil, errors.New("cannot find user with provided ID")
	}
	user = us.username2userMap[username]
	err = nil
	return
}

func (us *UserService) Add(user *models.User) {
	us.username2userMap[user.Username] = user
	us.usersDB.SaveUser(user)
}

func (us *UserService) Save(user *models.User) (stored bool) {
	return us.usersDB.SaveUser(user)
}

func (us *UserService) GetUserFromSpotify(accessToken string) (user *models.SpUser, err error) {
	body, err := getFromSpotify(config.Conf.SpotifyAPIURL, config.Conf.URLCurrentUser, accessToken)
	if err != nil {
		log.Printf(" >>> error getting current user playlists. details: %s\n", err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf(" >>> error getting current user playlists. details: %s\n", err.Error())
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetUserByRequestCookieID(r *http.Request) (user *models.User, err error) {
	cookieID, err := r.Cookie(constants.CookieUserIDKey)
	if err != nil {
		log.Printf(" >>> error, cannot find user by cookieID: %s\n", err.Error())
		return
	}

	user, err = us.GetUserByCookieID(cookieID.Value)
	if err != nil {
		log.Printf(" >>> %s\n", fmt.Sprintf(" >>> cannot find user by cookie. error: %s", err.Error()))
		return
	}
	return
}
