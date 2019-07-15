package api

import (
	"fmt"
	"io"
	"net/http"

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
		switch r.URL.Path {
		case "/api/ssfavtracks":
			handler.getFavTracksSnapshotsHandler(user.Username, false, w)
		case "/api/ssfavtracks/full":
			handler.getFavTracksSnapshotsHandler(user.Username, true, w)
		default:
			handler.getFavTracksSnapshot(user.Username, w, r)
		}
	case "DELETE":
		handler.deleteFavTracksSnapshots(user.Username, w, r)
	default:
		util.SendAPIErrorResp(w, "unknown/unsupported request method", http.StatusBadRequest)
	}
}

func (handler *FavTracksHandler) getFavTracksSnapshot(username string, w io.Writer, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Debugf(" > get fav tracks snapshot [%s]: username [%s]", timestamp, username)

	snapshot, err := services.UserPlaylist.GetFavTrakcsSnapshotByTimestamp(username, timestamp)
	if err != nil {
		log.Errorf(" >>> error while trying to get fav. tracks snapshot: %s", err.Error())
		util.SendAPIErrorResp(w, "Error occured: "+err.Error(), http.StatusNotFound)
		return
	}
	if snapshot == nil {
		log.Errorf(" >>> error while trying to get fav. tracks snapshot: snapshot is nil")
		util.SendAPIErrorResp(w, "Favorite tracks snapshot not found ", http.StatusNotFound)
		return
	}

	util.SendAPIOKRespWithData(w, "success", snapshot)
}

func (handler *FavTracksHandler) getFavTracksSnapshotsHandler(username string, loadAllData bool, w io.Writer) {
	log.WithFields(log.Fields{
		"loadAllData": loadAllData,
	}).Debugf(" > get fav tracks snapshots: username [%s]", username)

	sstracksRaw := services.UserPlaylist.GetAllFavTracksSnapshots(username)
	sstracks := []models.DTOFavTracksSnapshot{}
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
