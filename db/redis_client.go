package db

import (
	"encoding/json"
	"fmt"
	"log"

	c "github.com/2beens/spotilizer/constants"
	m "github.com/2beens/spotilizer/models"

	"gopkg.in/redis.v3"
)

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

	log.Printf(" > connected to redis %+v", options)
}

func StoreUserInfo(user *m.User) (stored bool) {
	auth, err := json.Marshal(user.Auth)
	if err != nil {
		fmt.Println(" >>> error while storing user info: " + err.Error())
		return false
	}
	authJSON := string(auth)
	log.Println(" > storing auth json:")
	log.Printf("%v\n", authJSON)
	rc.SAdd("users", fmt.Sprintf("%s::%s::%s", user.ID, user.Username, string(authJSON)))
	return true
}

func StoreCurrentTracks(user m.User) {
	log.Println(" > saving stored tracks ...")

}
