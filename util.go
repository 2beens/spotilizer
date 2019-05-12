package main

import (
	"math"
	"math/rand"
	"net/http"
	"time"
)

// generates a random string containing numbers and letters
func generateRandomString(length int) string {
	text := ""
	possible := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for i := 0; i < length; i++ {
		possibleLen := float64(len(possible))
		nextPossible := math.Floor(rand.Float64() * possibleLen)
		text += string(possible[int(nextPossible)])
		// text += possible.charAt(Math.floor(Math.random() * possible.length))
	}

	return text
}

// addCookie will apply a new cookie to the response of a http
// request, with the key/value this method is passed.
func addCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func cleearCookie(w http.ResponseWriter, name string) {
	expire := time.Now()
	cookie := http.Cookie{
		Name:    name,
		Value:   "",
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}
