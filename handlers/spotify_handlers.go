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

var clientID string
var clientSecret string

func SetCliendIdAndSecret(cID string, cSecret string) {
	clientID = cID
	clientSecret = cSecret
}

func GetSpotifyLoginHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		state := util.GenerateRandomString(16)
		util.AddCookie(w, c.CookieStateKey, state)

		q := url.Values{}
		q.Add("response_type", "code")
		q.Add("client_id", clientID)
		q.Add("scope", c.Permissions)
		q.Add("redirect_uri", fmt.Sprintf("%s/callback", serverURL))
		q.Add("state", state)

		redirectURL := "https://accounts.spotify.com/authorize?" + q.Encode()
		log.Println(" > /login, redirect to: " + redirectURL)
		http.Redirect(w, r, redirectURL, 302)
	}
}

func GetSaveCurrentTracksHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			// TOOD: redirect to error
			log.Printf(" >>> error while saving current user tracks: %v\n", err)
			http.Redirect(w, r, serverURL, 302)
			return
		}
		log.Printf(" > user ID: %s\n", userID.Value)
		if !s.Users.Exists(userID.Value) {
			// TOOD: redirect to error
			log.Printf(" >>> failed to find user, must login first\n")
			http.Redirect(w, r, serverURL, 302)
			return
		}

		user, _ := s.Users.Get(userID.Value)
		resp, err := s.UserPlaylist.GetSavedTracks(user.Auth)
		if err != nil {
			log.Printf(" >>> error while saving current user tracks: %v\n", err)
			return
		}
		log.Printf(" > tracks count: %d\n", len(resp))
		// TODO: return standardized resp message

		// TOOD: save tracks somewhere

		w.Write([]byte("track saved!"))
		return
	}
}

func GetSaveCurrentPlaylistsHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			// TOOD: redirect to error
			log.Printf(" >>> error while saving current user playlists: %v\n", err)
			return
		}

		log.Printf(" > user ID: %s\n", userID.Value)

		user, err := s.Users.Get(userID.Value)
		if err != nil {
			playlists, err := s.UserPlaylist.GetCurrentUserPlaylists(user.Auth)
			if err != nil {
				log.Printf(" >>> error while saving current user playlists: %v\n", err)
				return
			}
			log.Printf(" > playlists count: %d\n", len(playlists.Items))
			// TODO: return standardized resp message
			w.Write([]byte("playlists saved!"))
			return
		}
		log.Printf(" >>> failed to find user, must login first\n")
		http.Redirect(w, r, serverURL, 302)
	}
}

func GetRefreshTokenHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			log.Printf(" > refresh token failed, error: [%v]\n", err)
			// TODO: redirect to some error, or show error on the index page
			http.Redirect(w, r, serverURL, 302)
			return
		}
		user, err := s.Users.Get(userID.Value)
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

		var userID string
		user := s.Users.GetByUsername(spUser.ID)
		if user == nil {
			userID := util.GenerateRandomString(35)
			user = &m.User{Username: spUser.ID, ID: userID, Auth: authOptions}
			s.Users.Add(user)
		} else {
			userID = user.ID
		}

		util.AddCookie(w, c.CookieUserIDKey, userID)

		// redirect to index page with acces and refresh tokens
		util.RenderView(w, "index", m.ViewData{Message: "success", Error: "", Data: authOptions})
	}
}

func getAccessToken(data url.Values) m.SpotifyAuthOptions {
	body := postReq(data, "https://accounts.spotify.com", "/api/token/")
	authOptions := m.SpotifyAuthOptions{}
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
