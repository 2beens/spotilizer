package db

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	c "github.com/2beens/spotilizer/constants"
	m "github.com/2beens/spotilizer/models"

	"gopkg.in/redis.v3"
)

// redis is def not the best solution to persist the kind of data used in this server
// much better would be PostreSQL, or SQLite or so, but ... let's use redis for study reasons, but also not having to use SQL :)
var rc *redis.Client

func InitRedisClient(flashDB bool) {
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

	if flashDB {
		log.Println(" > will flush redis DB ...")
		FlushDB()
	}

	log.Printf(" > connected to redis %+v\n", options)
}

func FlushDB() {
	cmd := rc.FlushAll()
	res, err := cmd.Result()
	if err != nil {
		log.Printf(" >>> Flush DB error: %v\n", err)
		return
	}
	log.Printf(" > flush DB result: %s\n", res)
}

func SaveCookiesInfo(cookieID2usernameMap map[string]string) {
	log.Println(" > storing cookies data in DB ...")
	for id, username := range cookieID2usernameMap {
		log.Printf(" > [%s]: %s\n", id, username)
		idKey := "cookie::" + id
		cmd := rc.Set(idKey, username, 0)
		if err := cmd.Err(); err != nil {
			log.Printf(" >>> failed to store cookie ID for user: %s\n", username)
		}
	}
}

func GetCookiesInfo() (cookieID2usernameMap map[string]string) {
	cookieID2usernameMap = make(map[string]string)
	cmd := rc.Keys("cookie::*")
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get cookies info: %v\n", err)
		return nil
	}
	for _, cookieKey := range cmd.Val() {
		cookieID := strings.Split(cookieKey, "::")[1]
		cmd := rc.Get(cookieKey)
		if err := cmd.Err(); err != nil && err != redis.Nil {
			log.Printf(" >>> failed to get username for cookie ID %s: %v\n", cookieID, err)
		}
		username := cmd.Val()
		log.Printf(" > getting cookie from db [%s]: %s\n", cookieID, username)
		cookieID2usernameMap[cookieID] = username
	}
	return
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
	cmd := rc.Set(userKey, fmt.Sprintf("%s::%s", user.Username, authEncoded), 0)
	if err := cmd.Err(); err != nil {
		log.Printf(" >>> failed to store user info for user: %s\n", user.Username)
		return false
	}

	// TODO: save user tracks, playlists, and other data ...

	log.Printf(" > user [%s] saved to DB\n", user.Username)
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

func GetAllUsers() *[]m.User {
	cmd := rc.Keys("user::*")
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get all users: %s\n", err.Error())
		return nil
	}
	users := []m.User{}
	for _, userKey := range cmd.Val() {
		username := strings.Split(userKey, "::")[1]
		users = append(users, *GetUser(username))
	}
	return &users
}

func SavePlaylistsSnapshot(ps *m.PlaylistsSnapshot) (saved bool) {
	log.Printf(" > saving stored playlists [%d] for user [%s] ...\n", len(ps.Playlists), ps.Username)
	playlistsJSON, err := json.Marshal(ps.Playlists)
	if err != nil {
		log.Println(" >>> json marshalling error saving playlists to DB for user: " + ps.Username)
		return false
	}
	timestamp := strconv.FormatInt(ps.Timestamp.Unix(), 10)
	snapshotKey := fmt.Sprintf("playlistsshot::user::%s::timestamp::%s", ps.Username, timestamp)
	log.Printf(" > saving new playlist snapshot: [%s]\n", snapshotKey)
	cmd := rc.Set(snapshotKey, string(playlistsJSON), 0)
	if err := cmd.Err(); err != nil {
		log.Printf(" >>> failed to store playlists snapshot for user: %s\n", ps.Username)
		return false
	}
	log.Printf(" > user [%s] playlists snapshot saved to DB\n", ps.Username)
	return true
}

func GetPlaylistSnapshot(key string) *m.PlaylistsSnapshot {
	cmd := rc.Get(key)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get playlist snapshot [%s]: %s\n", key, err.Error())
		return nil
	}

	keyParts := strings.Split(key, "::")
	log.Printf(" > get playlists snapshot, key parts: %v\n", keyParts)
	username := keyParts[2]
	timestampStr := keyParts[4]
	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		log.Println(" >>> error while parsing playlist snapshot timestamp")
		return nil
	}
	timestamp := time.Unix(timestampInt, 0)

	playlistsJSON := cmd.Val()
	playlists := &[]m.SpPlaylist{}
	err = json.Unmarshal([]byte(playlistsJSON), playlists)
	if err != nil {
		log.Printf(" >>> failed to unmarshal playlists for [%s]: %s\n", key, err.Error())
		return nil
	}

	return &m.PlaylistsSnapshot{Username: username, Timestamp: timestamp, Playlists: *playlists}
}

func GetAllPlaylistsSnapshots(username string) *[]m.PlaylistsSnapshot {
	snapshotsKey := fmt.Sprintf("playlistsshot::user::%s::timestamp::*", username)
	cmd := rc.Keys(snapshotsKey)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get all playlists snapshots for user [%s]: %s\n", username, err.Error())
		return nil
	}
	plsnapshots := []m.PlaylistsSnapshot{}
	for _, skey := range cmd.Val() {
		ps := GetPlaylistSnapshot(skey)
		if ps != nil {
			plsnapshots = append(plsnapshots, *ps)
		}
	}
	return &plsnapshots
}
