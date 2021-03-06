package util

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/2beens/spotilizer/constants"
	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"

	log "github.com/sirupsen/logrus"
)

// GenerateRandomString generates a random string containing numbers and letters
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
func AddCookie(w *http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(*w, cookie)
}

func ClearCookie(w *http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   "",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(*w, cookie)
}

func GetUsernameByRequestCookieID(r *http.Request) (username string, found bool) {
	cookieID, err := r.Cookie(constants.CookieUserIDKey)
	if err != nil {
		// error ignored. can panic when invoked from incognito window
		return "", false
	}
	username, found = services.Users.GetUsernameByCookieID(cookieID.Value)
	return
}

func ReadSpotifyAuthData() (clientID string, clientSecret string, err error) {
	clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	log.Debug(" > client ID: " + clientID)
	log.Debug(" > client secret: " + clientSecret)
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
}

func SendAPIResp(w io.Writer, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf(" >>> Error while sending API response: %s\n", err.Error())
		return
	}
	_, err = w.Write(dataBytes)
	if err != nil {
		log.Println(" >>> error, failed to send API response. Writer error.")
	}
}

func SendAPIOKResp(w io.Writer, message string) {
	apiResp := models.APIResponse{Status: 200, Message: message}
	SendAPIResp(w, apiResp)
}

func SendAPIOKRespWithData(w io.Writer, message string, data interface{}) {
	apiResp := models.APIResponse{Status: 200, Message: message, Data: data}
	SendAPIResp(w, apiResp)
}

func SendAPIErrorResp(w io.Writer, message string, status int) {
	apiErr := models.SpAPIError{Error: models.SpError{Message: message, Status: status}}
	SendAPIResp(w, apiErr)
}
