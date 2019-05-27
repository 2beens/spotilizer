package handlers

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	c "github.com/2beens/spotilizer/constants"
	m "github.com/2beens/spotilizer/models"
	s "github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

func GetRefreshTokenHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookieID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			log.Printf(" > refresh token failed, error: [%v]\n", err)
			// TODO: redirect to some error, or show error on the index page
			http.Redirect(w, r, serverURL, 302)
			return
		}
		user, err := s.Users.Get(cookieID.Value)
		if err != nil {
			log.Println(" > refresh token failed, error: refresh_token param not found")
			// TODO: redirect to some error, or show error on the index page
			w.Write([]byte("missing refresh_token param"))
			return
		}
		log.Println(" > refresh token, value: " + user.Auth.RefreshToken)

		data := url.Values{}
		data.Set("refresh_token", user.Auth.RefreshToken)
		data.Set("grant_type", "refresh_token")
		newAuthOptions := getAccessToken(data)
		user.Auth = newAuthOptions

		// redirect to index page with acces and refresh tokens
		util.RenderView(w, "index", m.ViewData{Message: "success", Error: "", Data: user.Auth})
	}
}

func GetSpotifyCallbackHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		err, ok := q["error"]
		if ok {
			log.Printf(" > login failed, error: [%v]\n", err)
			// TODO: redirect to some error, or show error on the index page
			http.Redirect(w, r, serverURL, 302)
			return
		}

		code, codeOk := q["code"]
		state, stateOk := q["state"]
		storedStateCookie, sStateCookieErr := r.Cookie(c.CookieStateKey)
		if !codeOk || !stateOk {
			log.Println(" > login failed, error: some of the mandatory params not found")
			// TODO: redirect to some error, or show error on the index page
			http.Redirect(w, r, serverURL, 302)
			return
		}

		if storedStateCookie == nil || storedStateCookie.Value != state[0] || sStateCookieErr != nil {
			log.Printf(" > login failed, error: state cookie not found or state mismatch. more details [%v]\n", err)
			// TODO: redirect to some error, or show error on the index page
			http.Redirect(w, r, serverURL, 302)
			return
		}

		util.CleearCookie(w, c.CookieStateKey)
		data := url.Values{}
		data.Set("code", code[0])
		data.Set("redirect_uri", fmt.Sprintf("%s/callback", serverURL))
		data.Set("grant_type", "authorization_code")
		authOptions := getAccessToken(data)

		// get user info
		log.Println(" > getting user info from SP ...")
		spUser, userErr := s.Users.GetUserFromSpotify(authOptions)
		if userErr != nil {
			log.Println(" >>> error, cannot get user info from Spotify API.")
			// TODO: redirect to some error, or show error on the index page
			http.Redirect(w, r, serverURL, 302)
			return
		}
		log.Printf(" > gotten user [%s]\n", spUser.ID)

		var cookieID string
		user, _ := s.Users.Get(spUser.ID)
		if user == nil {
			cookieID = util.GenerateRandomString(45)
			user = &m.User{Username: spUser.ID, Auth: authOptions}
			s.Users.Add(user)
			s.Users.AddUserCookie(cookieID, user.Username)
			log.Printf(" > new user [%s] created and stored. cookie [%s]\n", user.Username, cookieID)
		} else {
			cID, cErr := s.Users.GetCookieIDByUsername(user.Username)
			if cErr != nil {
				cookieID = util.GenerateRandomString(45)
				log.Printf(" > generating and seding new cookie ID [%s] to client\n", cookieID)
				s.Users.AddUserCookie(cookieID, user.Username)
			} else {
				log.Println(" > using previous cookie ID: " + cookieID)
			}
			cookieID = cID
		}

		util.AddCookie(w, c.CookieUserIDKey, cookieID)

		// redirect to index page with acces and refresh tokens
		util.RenderView(w, "index", m.ViewData{Message: "success", Error: "", Data: authOptions})
	}
}

func getAccessToken(data url.Values) *m.SpotifyAuthOptions {
	body := postReq(data, "https://accounts.spotify.com", "/api/token/")
	authOptions := &m.SpotifyAuthOptions{}
	json.Unmarshal(body, &authOptions)
	return authOptions
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
		log.Printf(" >>> error making an auth post req: %v\n", err)
		// TODO: return error and check for it later
		return nil
	}
	defer resp.Body.Close()

	// redirect or return error if status != 200
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
