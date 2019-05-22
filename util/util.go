package util

import (
	"errors"
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	m "github.com/2beens/spotilizer/models"
)

// generates a random string containing numbers and letters
func GenerateRandomString(length int) string {
	text := ""
	possible := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for i := 0; i < length; i++ {
		possibleLen := float64(len(possible))
		nextPossible := math.Floor(rand.Float64() * possibleLen)
		text += string(possible[int(nextPossible)])
	}

	return text
}

// AddCookie will apply a new cookie to the response of a http
// request, with the key/value this method is passed.
func AddCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func CleearCookie(w http.ResponseWriter, name string) {
	expire := time.Now()
	cookie := http.Cookie{
		Name:    name,
		Value:   "",
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func ReadSpotifyAuthData() (clientID string, clientSecret string, err error) {
	clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	log.Println(" > client ID: " + clientID)
	log.Println(" > client secret: " + clientSecret)
	if clientID == "" {
		return "", "", errors.New(" >>> error, client ID missing. set it using env [SPOTIFY_CLIENT_ID]")
	}
	if clientSecret == "" {
		return "", "", errors.New(" >>> error, client secret missing. set it using env [SPOTIFY_CLIENT_SECRET]")
	}
	return
}

func LoggingSetup(logFileName string) {
	if logFileName == "" {
		log.SetOutput(os.Stdout)
		return
	}

	if !strings.HasSuffix(logFileName, ".log") {
		logFileName += ".log"
	}

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Panicf("failed to open log file %q: %s", logFileName, err)
	}

	log.SetOutput(logFile)
	log.SetFlags(5)
}

// templates cheatsheet
// https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet
func RenderView(w http.ResponseWriter, page string, viewData m.ViewData) {
	// TODO: parse the template once and reuse it
	files := []string{
		"public/views/layouts/layout.html",
		"public/views/layouts/footer.html",
		"public/views/layouts/navbar.html",
		"public/views/" + page + ".html",
	}
	t, err := template.New("layout").ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "layout", viewData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}