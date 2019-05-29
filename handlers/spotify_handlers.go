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
			log.Printf(" >>> %s\n", fmt.Sprintf(" >>> cookie error while saving current user tracks: %s", err.Error()))
			util.SendAPIErrorResp(w, "Not available when logged off", 400)
			return
		}

		user, err := s.Users.GetUserByCookieID(cookieID.Value)
		if err != nil {
			log.Printf(" >>> %s\n", fmt.Sprintf(" >>> user/cookie error while saving current user tracks: %s", err.Error()))
			util.SendAPIErrorResp(w, "Not available when logged off", 400)
			return
		}

		log.Printf(" > get fav tracks: cookie [%s], username [%s]\n", cookieID.Value, user.Username)

		tracks, apiErr := s.UserPlaylist.GetSavedTracks(user.Auth)
		if apiErr != nil {
			log.Printf(" >>> error while saving current user tracks: %v\n", apiErr)
			util.SendAPIErrorResp(w, apiErr.Error.Message, apiErr.Error.Status)
			return
		}

		log.Printf(" > tracks count: %d\n", len(tracks))
		user.FavTracks = &tracks

		// TODO: return standardized resp message

		// TOOD: save tracks somewhere (redis)

		util.SendAPIOKResp(w, fmt.Sprintf("%d tracks saved successfully", len(tracks)))
		return
	}
}

func GetSaveCurrentPlaylistsHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookieID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			log.Printf(" >>> %s\n", fmt.Sprintf(" >>> cookie error while saving current user playlists: %s", err.Error()))
			util.SendAPIErrorResp(w, "Not available when logged off", 400)
			return
		}

		user, err := s.Users.GetUserByCookieID(cookieID.Value)
		if err != nil {
			log.Printf(" >>> %s\n", fmt.Sprintf(" >>> user/cookie error while saving current user playlists: %s", err.Error()))
			util.SendAPIErrorResp(w, "Not available when logged off", 400)
			return
		}

		log.Printf(" > user ID: %s\n", cookieID.Value)

		playlists, err := s.UserPlaylist.GetCurrentUserPlaylists(user.Auth)
		if err != nil {
			log.Printf(" >>> error while saving current user playlists: %v\n", err)
			util.SendAPIErrorResp(w, err.Error(), 400)
			return
		}
		log.Printf(" > playlists count: %d\n", len(playlists.Items))

		// TODO: save playlists

		util.SendAPIOKResp(w, fmt.Sprintf("%d playlists saved successfully", len(playlists.Items)))
		return
	}
}
