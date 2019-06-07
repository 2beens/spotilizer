package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/2beens/spotilizer/models"
	m "github.com/2beens/spotilizer/models"
	s "github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

func SaveCurrentTracksHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.Users.GetUserByRequestCookieID(r)
	if err != nil {
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Printf(" > save fav tracks: username [%s]\n", user.Username)

	tracks, apiErr := s.UserPlaylist.DownloadSavedFavTracks(user.Auth)
	if apiErr != nil {
		log.Printf(" >>> error while saving current user tracks: %v\n", apiErr)
		util.SendAPIErrorResp(w, apiErr.Error.Message, apiErr.Error.Status)
		return
	}

	log.Printf(" > tracks count: %d\n", len(tracks))

	// save tracks to DB
	tracksSnapshot := &m.FavTracksSnapshot{Username: user.Username, Timestamp: time.Now(), Tracks: tracks}
	saved := s.UserPlaylist.SaveFavTracksSnapshot(tracksSnapshot)
	if saved {
		util.SendAPIOKResp(w, fmt.Sprintf("%d favorite tracks saved successfully", len(tracks)))
	} else {
		util.SendAPIErrorResp(w, "Favorite tracks not saved. Server internal error.", http.StatusInternalServerError)
	}
}

func SaveCurrentPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.Users.GetUserByRequestCookieID(r)
	if err != nil {
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Printf(" > save playlists: username: %s\n", user.Username)

	playlists, apiErr := s.UserPlaylist.DownloadCurrentUserPlaylists(user.Auth)
	if apiErr != nil {
		log.Printf(" >>> error while saving current user playlists: %v\n", apiErr)
		util.SendAPIErrorResp(w, apiErr.Error.Message, apiErr.Error.Status)
		return
	}

	log.Printf(" > playlists count: %d\n", len(playlists))

	snapshotPlaylists := []models.PlaylistSnapshot{}
	for _, pl := range playlists {
		playlistTracks, apiErr := s.UserPlaylist.DownloadPlaylistTracks(user.Auth, pl.Tracks.Href, pl.Tracks.Total)
		if apiErr != nil {
			log.Printf(" >>> error while saving current user playlists: %v\n", apiErr)
			playlistTracks = []models.SpPlaylistTrack{}
		}
		log.Printf(" > received [%d] tracks for playlist [%s]\n", len(playlistTracks), pl.Name)
		plSnapshot := models.PlaylistSnapshot{
			Playlist: pl,
			Tracks:   playlistTracks,
		}
		snapshotPlaylists = append(snapshotPlaylists, plSnapshot)
	}

	// save playlists to DB
	playlistsSnapshot := &m.PlaylistsSnapshot{Username: user.Username, Timestamp: time.Now(), Playlists: snapshotPlaylists}
	saved := s.UserPlaylist.SavePlaylistsSnapshot(playlistsSnapshot)
	if saved {
		util.SendAPIOKResp(w, fmt.Sprintf("%d playlists saved successfully", len(playlists)))
	} else {
		util.SendAPIErrorResp(w, "Playlists not saved. Server internal error.", http.StatusInternalServerError)
	}
}
