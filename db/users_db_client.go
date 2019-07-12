package db

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	b64 "encoding/base64"

	m "github.com/2beens/spotilizer/models"
	"gopkg.in/redis.v3"
)

type UsersDBClient interface {
	SaveUser(user *m.User) (stored bool)
	GetUser(username string) *m.User
	GetAllUsers() *[]m.User
}

type UsersDB struct{}

func (uDB UsersDB) SaveUser(user *m.User) (stored bool) {
	auth, err := json.Marshal(user.Auth)
	if err != nil {
		fmt.Println(" >>> error while storing user info: " + err.Error())
		return false
	}
	authJSON := string(auth)
	authEncoded := b64.StdEncoding.EncodeToString([]byte(authJSON))
	userKey := "user::" + user.Username
	cmd := rc.Set(userKey, fmt.Sprintf("%s::%s", user.Username, authEncoded), 0)
	if err := cmd.Err(); err != nil {
		log.Printf(" >>> failed to store user info for user: %s\n", user.Username)
		return false
	}

	// TODO: save user tracks, playlists, and other data ... ??

	log.Printf(" > user [%s] saved to DB\n", user.Username)
	return true
}

// GetUser returns a user object from storage (redis) by username
func (uDB UsersDB) GetUser(username string) *m.User {
	cmd := rc.Get("user::" + username)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get user %s: %v\n", username, err)
		return nil
	}
	userStringData := cmd.Val()
	userData := strings.Split(userStringData, "::")
	authDecoded, err := b64.StdEncoding.DecodeString(userData[1])
	if err != nil {
		log.Printf(" >>> failed to get user %s: %v\n", username, err)
		return nil
	}
	auth := &m.SpotifyAuthOptions{}
	err = json.Unmarshal(authDecoded, auth)
	if err != nil {
		log.Printf(" >>> failed to get user %s: %v\n", username, err)
		return nil
	}
	return &m.User{Username: username, Auth: auth}
}

func (uDB UsersDB) GetAllUsers() *[]m.User {
	cmd := rc.Keys("user::*")
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get all users: %s\n", err.Error())
		return nil
	}
	users := []m.User{}
	for _, userKey := range cmd.Val() {
		username := strings.Split(userKey, "::")[1]
		users = append(users, *uDB.GetUser(username))
	}
	return &users
}
