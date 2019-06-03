package db

import (
	"log"
	"strings"

	"gopkg.in/redis.v3"
)

type CookiesDBClient interface {
	SaveCookiesInfo(cookieID2usernameMap map[string]string)
	GetCookiesInfo() (cookieID2usernameMap map[string]string)
}

type CookiesDB struct{}

func (cDB CookiesDB) SaveCookiesInfo(cookieID2usernameMap map[string]string) {
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

func (cDB CookiesDB) GetCookiesInfo() (cookieID2usernameMap map[string]string) {
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
