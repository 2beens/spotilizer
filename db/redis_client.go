package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/2beens/spotilizer/constants"

	"gopkg.in/redis.v3"
)

// redis is def not the best solution to persist the kind of data used in this server
// much better would be PostreSQL, or SQLite or so, but ... let's use redis for study reasons, but also not having to use SQL :)
var rc *redis.Client

// Cookies client for all cookies related data
var cookiesDBClient CookiesDBClient
var usersDBClient UsersDBClient
var spotifyDBClient SpotifyDBClient

func InitRedisClient(flashDB bool) {
	log.Println(" > initializing redis ...")
	options := &redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%s", constants.IPAddress, constants.RedisPort), // localhost:6379
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

	cookiesDBClient = &CookiesDB{}
	usersDBClient = &UsersDBRedisClient{}
	spotifyDBClient = &SpotifyDB{}

	log.Printf(" > connected to redis %+v\n", options)
}

func GetCookiesDBClient() CookiesDBClient {
	return cookiesDBClient
}

func GetUsersDBClient() UsersDBClient {
	return usersDBClient
}

func GetSpotifyDBClient() SpotifyDBClient {
	return spotifyDBClient
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
