package handlers

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	c "github.com/2beens/spotilizer/constants"
	m "github.com/2beens/spotilizer/models"
	s "github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

var clientID string
var clientSecret string

func SetCliendIDAndSecret(cID string, cSecret string) {
	clientID = cID
	clientSecret = cSecret
}

func GetSpotifyLoginHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		state := util.GenerateRandomString(16)
		util.AddCookie(&w, c.CookieStateKey, state)

		q := url.Values{}
		q.Add("response_type", "code")
		q.Add("client_id", clientID)
		q.Add("scope", c.Permissions)
		q.Add("redirect_uri", fmt.Sprintf("%s/callback", serverURL))
		q.Add("state", state)

		redirectURL := "https://accounts.spotify.com/authorize?" + q.Encode()
		log.Trace(" > /login, redirect to: " + redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookieID, err := r.Cookie(c.CookieUserIDKey)
	if err != nil {
		log.Debugf(" > refresh token failed, cannot find user by cookie ID, error: [%s]", err.Error())
		util.SendAPIErrorResp(w, "Cannot find user by cookie, refresh token failed", 400)
		return
	}
	user, err := s.Users.GetUserByCookieID(cookieID.Value)
	if err != nil {
		log.Debugf(" > refresh token failed, cannot find user by cookie ID [%s]", cookieID.Value)
		util.SendAPIErrorResp(w, "Cannot find user by cookie, refresh token failed", 400)
		return
	}
	log.Traceln(" > refresh token, value: " + user.Auth.RefreshToken)

	data := url.Values{}
	data.Set("refresh_token", user.Auth.RefreshToken)
	data.Set("grant_type", "refresh_token")
	newAuthOptions, authErr := getAccessToken(data)
	if authErr != nil {
		util.SendAPIErrorResp(w, "Internal server error during refresh token. Try again later.", http.StatusInternalServerError)
		return
	}
	user.Auth = newAuthOptions

	s.Users.Save(user)

	log.Tracef(" > refresh token success for user [%s]", user.Username)
	util.SendAPIOKResp(w, "success")
}

func GetSpotifyCallbackHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		err, ok := q["error"]
		if ok {
			log.Infof(" > login failed, error: [%v]", err)
			util.RenderView(w, "error", m.ErrorViewData{Title: "Spotify Login",
				Error: "Login to Spotify failed: " + strings.Join(err, ", ")})
			return
		}

		code, codeOk := q["code"]
		state, stateOk := q["state"]
		storedStateCookie, sStateCookieErr := r.Cookie(c.CookieStateKey)
		if !codeOk || !stateOk {
			log.Debug(" > login failed, error: some of the mandatory params not found")
			util.RenderView(w, "error", m.ErrorViewData{Title: "Spotify Login",
				Error: "Login to Spotify failed: login failed, error: some of the mandatory params not found"})
			return
		}

		if storedStateCookie == nil || storedStateCookie.Value != state[0] || sStateCookieErr != nil {
			log.Debugf(" > login failed, error: state cookie not found or state mismatch. more details [%v]", err)
			util.RenderView(w, "error", m.ErrorViewData{Title: "Spotify Login",
				Error: "Login to Spotify failed: state cookie not found or state mismatch"})
			return
		}

		util.CleearCookie(&w, c.CookieStateKey)
		data := url.Values{}
		data.Set("code", code[0])
		data.Set("redirect_uri", fmt.Sprintf("%s/callback", serverURL))
		data.Set("grant_type", "authorization_code")
		authOptions, authErr := getAccessToken(data)
		if authErr != nil {
			util.RenderErrorView(w, "", "Login Failed", http.StatusInternalServerError, "Internal server error during login. Try again later.")
			return
		}

		// get user info
		log.Debugln(" > getting user info from SP ...")
		spUser, userErr := s.Users.GetUserFromSpotify(authOptions.AccessToken)
		if userErr != nil {
			log.Debugf(" >>> error, cannot get user info from Spotify API.")
			util.RenderView(w, "error", m.ErrorViewData{Title: "Spotify Login",
				Error: "Login to Spotify failed: error, cannot get user info from Spotify API"})
			return
		}
		log.Debugf(" > gotten user [%s]\n", spUser.ID)

		var cookieID string
		user, _ := s.Users.Get(spUser.ID)
		if user == nil {
			cookieID = util.GenerateRandomString(45)
			user = &m.User{Username: spUser.ID, Auth: authOptions}
			s.Users.Add(user)
			s.Users.AddUserCookie(cookieID, user.Username)
			log.Debugf(" > new user [%s] created and stored. cookie [%s]", user.Username, cookieID)
		} else {
			cID, cErr := s.Users.GetCookieIDByUsername(user.Username)
			if cErr != nil {
				cookieID = util.GenerateRandomString(45)
				log.Debugf(" > generating and seding new cookie ID [%s] to client", cookieID)
				s.Users.AddUserCookie(cookieID, user.Username)
			} else {
				log.Debug(" > using previous cookie ID: " + cID)
			}
			cookieID = cID
			user.Auth = authOptions
			s.Users.Save(user)
		}

		util.AddCookie(&w, c.CookieUserIDKey, cookieID)

		GetIndexHandler(user.Username)(w, r)
	}
}

func getAccessToken(data url.Values) (auth *m.SpotifyAuthOptions, err error) {
	body := postReq(data, "https://accounts.spotify.com", "/api/token/")
	auth = &m.SpotifyAuthOptions{}
	err = json.Unmarshal(body, &auth)
	if err != nil {
		log.Warn(" >>> error while unmarshaling auth options: " + err.Error())
		auth = nil
	}
	return
}

func postReq(data url.Values, uri string, path string) []byte {
	u, _ := url.ParseRequestURI(uri)
	u.Path = path

	client := &http.Client{}
	r, _ := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode())) // URL-encoded payload
	authEncoding := b64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	r.Header.Add("Authorization", "Basic "+authEncoding)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		log.Infof(" >>> error making an auth post req: %s", err.Error())
		// TODO: return error and check for it later
		return nil
	}
	defer resp.Body.Close()

	// redirect or return error if status != 200
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
