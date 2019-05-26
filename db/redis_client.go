package db

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	c "github.com/2beens/spotilizer/constants"
	m "github.com/2beens/spotilizer/models"

	"gopkg.in/redis.v3"
)

// redis is def not the best solution to persist the kind of data used in this server
// much better would be PostreSQL, or SQLite or so, but ... let's use redis for study reasons, but also not having to use SQL :)
var rc *redis.Client

func InitRedisClient() {
	log.Println(" > initializing redis ...")
	options := &redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%s", c.IPAddress, c.RedisPort), // localhost:6379
		DB:      int64(6),
	}

	rc = redis.NewClient(options)

	if err := rc.Ping().Err(); err != nil {
		log.Panicf(" >>> failed to ping redis %+v", options)
	}

	log.Printf(" > connected to redis %+v\n", options)
}

func SaveUser(user *m.User) (stored bool) {
	auth, err := json.Marshal(user.Auth)
	if err != nil {
		fmt.Println(" >>> error while storing user info: " + err.Error())
		return false
	}
	authJSON := string(auth)
	authEncoded := b64.StdEncoding.EncodeToString([]byte(authJSON))
	userKey := "user::" + user.Username
	cmd := rc.Set(userKey, fmt.Sprintf("%s::%s::%s", user.ID, user.Username, authEncoded), 0)
	if err := cmd.Err(); err != nil {
		log.Printf(" >>> failed to store user info for user: %s, [%s]\n", user.Username, user.ID)
		return false
	}
	return true
}

// GetUser returns a user object from storage (redis) by username
func GetUser(username string) *m.User {
	cmd := rc.Get("user::" + username)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get user %s: %v\n", username, err)
		return nil
	}
	userStringData := cmd.Val()
	userData := strings.Split(userStringData, "::")
	authDecoded, err := b64.StdEncoding.DecodeString(userData[2])
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
	return &m.User{Username: username, ID: userData[0], Auth: auth}
}

func GetAllUsers() *[]m.User {
	cmd := rc.Keys("user::*")
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get all users: %v\n", err)
		return nil
	}
	users := []m.User{}
	for _, userKey := range cmd.Val() {
		username := strings.Split(userKey, "::")[1]
		users = append(users, *GetUser(username))
	}
	return &users
}

func StoreCurrentTracks(user m.User) {
	log.Println(" > saving stored tracks ...")

}
