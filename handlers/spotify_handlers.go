package handlers

import (
	"fmt"
	"log"
	"net/http"

	c "github.com/2beens/spotilizer/constants"
	s "github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

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

		playlists, apiErr := s.UserPlaylist.GetCurrentUserPlaylists(user.Auth)
		if apiErr != nil {
			log.Printf(" >>> error while saving current user playlists: %v\n", err)
			util.SendAPIErrorResp(w, apiErr.Error.Message, apiErr.Error.Status)
			return
		}

		log.Printf(" > playlists count: %d\n", len(playlists))
		user.Playlists = &playlists

		// TODO: save playlists

		util.SendAPIOKResp(w, fmt.Sprintf("%d playlists saved successfully", len(playlists)))
		return
	}
}
