package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	c "github.com/2beens/spotilizer/constants"
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
		util.AddCookie(&w, c.CookieStateKey, state)

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
		cookieID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			// TOOD: redirect to error
			log.Printf(" >>> error while saving current user tracks: %v\n", err)
			http.Redirect(w, r, serverURL, 302)
			return
		}

		log.Println(" > user cookie: " + cookieID.Value)

		user, err := s.Users.GetUserByCookieID(cookieID.Value)
		if err != nil {
			// TOOD: redirect to error
			log.Printf(" >>> failed to find user, must login first\n")
			http.Redirect(w, r, serverURL, 302)
			return
		}

		tracks, err := s.UserPlaylist.GetSavedTracks(user.Auth)
		if err != nil {
			log.Printf(" >>> error while saving current user tracks: %v\n", err)
			return
		}

		log.Printf(" > tracks count: %d\n", len(tracks))
		user.FavTracks = &tracks

		// TODO: return standardized resp message

		// TOOD: save tracks somewhere (redis)

		w.Write([]byte("track saved!"))
		return
	}
}

func GetSaveCurrentPlaylistsHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookieID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			// TOOD: redirect to error
			log.Printf(" >>> error while saving current user playlists: %v\n", err)
			return
		}

		log.Printf(" > user ID: %s\n", cookieID.Value)

		user, err := s.Users.Get(cookieID.Value)
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
