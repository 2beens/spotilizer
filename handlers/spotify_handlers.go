package handlers

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

func SaveCurrentTracksHandler(w http.ResponseWriter, r *http.Request) {
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Debugf(" > save fav tracks: username [%s]", user.Username)

	tracks, apiErr := services.UserPlaylist.DownloadSavedFavTracks(user.Auth.AccessToken)
	if apiErr != nil {
		log.Infof(" >>> error while saving current user tracks: %v", apiErr)
		util.SendAPIErrorResp(w, apiErr.Error.Message, apiErr.Error.Status)
		return
	}

	log.Tracef(" > tracks count: %d", len(tracks))

	// save tracks to DB
	tracksSnapshot := &models.FavTracksSnapshot{Username: user.Username, Timestamp: time.Now(), Tracks: tracks}
	saved := services.UserPlaylist.SaveFavTracksSnapshot(tracksSnapshot)
	if saved {
		util.SendAPIOKResp(w, fmt.Sprintf("%d favorite tracks saved successfully", len(tracks)))
	} else {
		util.SendAPIErrorResp(w, "Favorite tracks not saved. Server internal error.", http.StatusInternalServerError)
	}
}

func SaveCurrentPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Debugf(" > save playlists: username: %s", user.Username)

	playlists, apiErr := services.UserPlaylist.DownloadCurrentUserPlaylists(user.Auth.AccessToken)
	if apiErr != nil {
		log.Infof(" >>> error while saving current user playlists: %v", apiErr)
		util.SendAPIErrorResp(w, apiErr.Error.Message, apiErr.Error.Status)
		return
	}

	log.Tracef(" > playlists count: %d", len(playlists))

	snapshotPlaylists := []models.PlaylistSnapshot{}
	for _, pl := range playlists {
		playlistTracks, apiErr := services.UserPlaylist.DownloadPlaylistTracks(user.Auth.AccessToken, pl.Tracks.Href, pl.Tracks.Total)
		if apiErr != nil {
			log.Warnf(" >>> error while saving current user playlists: %v", apiErr)
			playlistTracks = []models.SpPlaylistTrack{}
		}
		log.Tracef(" > received [%d] tracks for playlist [%s]", len(playlistTracks), pl.Name)
		plSnapshot := models.PlaylistSnapshot{
			Playlist: pl,
			Tracks:   playlistTracks,
		}
		snapshotPlaylists = append(snapshotPlaylists, plSnapshot)
	}

	// save playlists to DB
	playlistsSnapshot := &models.PlaylistsSnapshot{Username: user.Username, Timestamp: time.Now(), Playlists: snapshotPlaylists}
	saved := services.UserPlaylist.SavePlaylistsSnapshot(playlistsSnapshot)
	if saved {
		util.SendAPIOKResp(w, fmt.Sprintf("%d playlists saved successfully", len(playlists)))
	} else {
		util.SendAPIErrorResp(w, "Playlists not saved. Server internal error.", http.StatusInternalServerError)
	}
}
