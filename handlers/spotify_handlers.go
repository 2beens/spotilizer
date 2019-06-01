package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/2beens/spotilizer/db"

	c "github.com/2beens/spotilizer/constants"
	m "github.com/2beens/spotilizer/models"
	s "github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

func SaveCurrentTracksHandler(w http.ResponseWriter, r *http.Request) {
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

	// save tracks to DB
	tracksSnapshot := &m.FavTracksSnapshot{Username: user.Username, Timestamp: time.Now(), Tracks: tracks}
	saved := db.SaveFavTracksSnapshot(tracksSnapshot)
	if saved {
		util.SendAPIOKResp(w, fmt.Sprintf("%d favorite tracks saved successfully", len(tracks)))
	} else {
		util.SendAPIErrorResp(w, "Favorite tracks not saved. Server internal error.", 500)
	}
	return
}

func SaveCurrentPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Printf(" >>> error while saving current user playlists: %v\n", apiErr)
		util.SendAPIErrorResp(w, apiErr.Error.Message, apiErr.Error.Status)
		return
	}

	log.Printf(" > playlists count: %d\n", len(playlists))
	user.Playlists = &playlists

	// save playlists to DB
	playlistsSnapshot := &m.PlaylistsSnapshot{Username: user.Username, Timestamp: time.Now(), Playlists: playlists}
	saved := db.SavePlaylistsSnapshot(playlistsSnapshot)
	if saved {
		util.SendAPIOKResp(w, fmt.Sprintf("%d playlists saved successfully", len(playlists)))
	} else {
		util.SendAPIErrorResp(w, "Playlists not saved. Server internal error.", 500)
	}
	return
}
