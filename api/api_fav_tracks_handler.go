package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type FavTracksHandler struct{}

func NewFavTracksHandler() *FavTracksHandler {
	return &FavTracksHandler{}
}

func (handler *FavTracksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		log.Errorf(" >>> API fav. tracks handler: user/cookie error: %s", err.Error())
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	switch r.Method {
	case "GET":
		if r.URL.Path == "/api/ssfavtracks" {
			handler.getFavTracksSnapshots(user.Username, false, w)
		} else if r.URL.Path == "/api/ssfavtracks/full" {
			handler.getFavTracksSnapshots(user.Username, true, w)
		} else if strings.HasPrefix(r.URL.Path, "/api/ssfavtracks/diff/") {
			handler.getFavTracksDiff(user, w, r)
		} else if r.URL.Path == "/api/ssfavtracks/{timestamp}" {
			handler.getFavTracksSnapshot(user.Username, w, r)
		} else {
			util.SendAPIErrorResp(w, "unknown path", http.StatusBadRequest)
		}
	case "DELETE":
		handler.deleteFavTracksSnapshots(user.Username, w, r)
	default:
		util.SendAPIErrorResp(w, "unknown/unsupported request method", http.StatusBadRequest)
	}
}

func (handler *FavTracksHandler) getFavTracksDiff(user *models.User, w io.Writer, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Debugf(" > get fav tracks diff for snapshot [%s]: username [%s]", timestamp, user.Username)

	snapshot, err := services.UserPlaylist.GetFavTracksSnapshotByTimestamp(user.Username, timestamp)
	if err != nil {
		log.Errorf(" >>> error while trying to get fav. tracks snapshot: %s", err.Error())
		util.SendAPIErrorResp(w, "Error occurred: "+err.Error(), http.StatusNotFound)
		return
	}
	if snapshot == nil {
		log.Errorf(" >>> error while trying to get fav. tracks snapshot: snapshot is nil")
		util.SendAPIErrorResp(w, "Favorite tracks snapshot not found ", http.StatusNotFound)
		return
	}

	// now get the current fav tracks, and make a diff relative to "snapshot" object
	currentTracks, apiErr := services.UserPlaylist.DownloadSavedFavTracks(user.Auth.AccessToken)
	if apiErr != nil {
		log.Infof(" >>> error while getting current tracks diff: %v", apiErr)
		util.SendAPIErrorResp(w, apiErr.Error.Message, apiErr.Error.Status)
		return
	}

	var newTracks []models.SpAddedTrack
	var removedTracks []models.SpAddedTrack
	for _, t := range currentTracks {
		if !containsTrack(t, snapshot.Tracks) {
			newTracks = append(newTracks, t)
		}
	}
	for _, t := range snapshot.Tracks {
		if !containsTrack(t, currentTracks) {
			removedTracks = append(removedTracks, t)
		}
	}

	log.Debugf(" > fav tracks [%s] diff. found [%d] new tracks and [%d] removed tracks", timestamp, len(newTracks), len(removedTracks))

	util.SendAPIOKRespWithData(w, "success", struct {
		NewTracks     []models.SpAddedTrack `json:"newTracks"`
		RemovedTracks []models.SpAddedTrack `json:"removedTracks"`
	}{
		newTracks,
		removedTracks,
	})
}

func containsTrack(track models.SpAddedTrack, tracks []models.SpAddedTrack) bool {
	for _, t := range tracks {
		if track.Track.Album.ID != t.Track.Album.ID {
			continue
		}
		if track.Track.ID != t.Track.ID {
			continue
		}
		return true
	}
	return false
}

func (handler *FavTracksHandler) getFavTracksSnapshot(username string, w io.Writer, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Debugf(" > get fav tracks snapshot [%s]: username [%s]", timestamp, username)

	snapshot, err := services.UserPlaylist.GetFavTracksSnapshotByTimestamp(username, timestamp)
	if err != nil {
		log.Errorf(" >>> error while trying to get fav. tracks snapshot: %s", err.Error())
		util.SendAPIErrorResp(w, "Error occurred: "+err.Error(), http.StatusNotFound)
		return
	}
	if snapshot == nil {
		log.Errorf(" >>> error while trying to get fav. tracks snapshot: snapshot is nil")
		util.SendAPIErrorResp(w, "Favorite tracks snapshot not found ", http.StatusNotFound)
		return
	}

	util.SendAPIOKRespWithData(w, "success", snapshot)
}

func (handler *FavTracksHandler) getFavTracksSnapshots(username string, loadAllData bool, w io.Writer) {
	log.WithFields(log.Fields{
		"loadAllData": loadAllData,
	}).Debugf(" > get fav tracks snapshots: username [%s]", username)

	sstracksRaw := services.UserPlaylist.GetAllFavTracksSnapshots(username)
	var sstracks []models.DTOFavTracksSnapshot
	for _, tracksssRaw := range sstracksRaw {
		tracksss := models.DTOFavTracksSnapshot{
			Timestamp:   tracksssRaw.Timestamp.Unix(),
			TracksCount: len(tracksssRaw.Tracks),
			Tracks:      []models.DTOTrack{},
		}
		if loadAllData {
			for _, trRaw := range tracksssRaw.Tracks {
				tracksss.Tracks = append(tracksss.Tracks, models.SpAddedTrack2dtoTrack(trRaw))
			}
		}
		sstracks = append(sstracks, tracksss)
	}

	util.SendAPIOKRespWithData(w, "success", sstracks)
}

func (handler *FavTracksHandler) deleteFavTracksSnapshots(username string, w io.Writer, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Debugf(" > delete fav tracks snapshot [%s]: username [%s]", timestamp, username)

	snapshot, err := services.UserPlaylist.DeleteFavTracksSnapshot(username, timestamp)
	if err != nil {
		log.Errorf(" >>> error while trying to delete fav. tracks snapshot: %s", err.Error())
		util.SendAPIErrorResp(w, "Error occured: "+err.Error(), http.StatusNotFound)
		return
	}
	if snapshot == nil {
		log.Errorf(" >>> error while trying to delete fav. tracks snapshot: snapshot is nil")
		util.SendAPIErrorResp(w, "Favorite tracks snapshot not deleted: not found ", http.StatusNotFound)
		return
	}

	util.SendAPIOKResp(w, fmt.Sprintf("Favorite tracks snapshot [%s] successfully deleted.", snapshot.Timestamp))
}
