package api

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
	"github.com/gorilla/mux"
)

type PlaylistsHandler struct{}

func NewPlaylistsHandler() *PlaylistsHandler {
	return &PlaylistsHandler{}
}

func (handler *PlaylistsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		log.Errorf(" >>> API playlists handler: user/cookie error: %s", err.Error())
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	switch r.Method {
	case "GET":
		switch r.URL.Path {
		case "/api/ssplaylists":
			handler.getPlaylistsSnapshotsHandler(user.Username, false, w, r)
		default:
			handler.getPlaylistsSnapshotsHandler(user.Username, true, w, r)
		}
	case "DELETE":
		handler.deletePlaylistsSnapshot(user.Username, w, r)
	default:
		util.SendAPIErrorResp(w, "unknown/unsupported request method", http.StatusBadRequest)
	}
}

func (handler *PlaylistsHandler) deletePlaylistsSnapshot(username string, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Debugf(" > delete playlists snapshot [%s]: username [%s]", timestamp, username)

	snapshot, err := services.UserPlaylist.DeletePlaylistsSnapshot(username, timestamp)
	if err != nil {
		log.Errorf(" >>> error while trying to delete playlists snapshot: %s", err.Error())
		util.SendAPIErrorResp(w, "Error occured: "+err.Error(), http.StatusNotFound)
		return
	}
	if snapshot == nil {
		log.Errorf(" >>> error while trying to delete playlists snapshot: snapshot is nil")
		util.SendAPIErrorResp(w, "Playlists snapshot not deleted: not found ", http.StatusNotFound)
		return
	}

	util.SendAPIOKResp(w, fmt.Sprintf("Playlists snapshot [%s] successfully deleted.", snapshot.Timestamp))
}

func (handler *PlaylistsHandler) getPlaylistsSnapshotsHandler(username string, loadAllData bool, w http.ResponseWriter, r *http.Request) {
	log.Debugf(" > get playlists snapshots: username [%s]", username)
	ssplaylists := handler.preparePlaylistsSnapshots(username, loadAllData)
	util.SendAPIOKRespWithData(w, "success", ssplaylists)
}

func (handler *PlaylistsHandler) preparePlaylistsSnapshots(username string, loadTracks bool) []models.DTOPlaylistSnapshot {
	ssplaylistsRaw := services.UserPlaylist.GetAllPlaylistsSnapshots(username)
	ssplaylists := []models.DTOPlaylistSnapshot{}
	for _, plssRaw := range ssplaylistsRaw {
		plss := models.DTOPlaylistSnapshot{
			Timestamp: plssRaw.Timestamp.Unix(),
			Playlists: []models.DTOPlaylist{},
		}
		rawTracks := []models.SpPlaylistTrack{}
		for _, plRaw := range plssRaw.Playlists {
			if loadTracks {
				rawTracks = plRaw.Tracks
			}
			plss.Playlists = append(plss.Playlists, models.SpPlaylist2dtoPlaylist(plRaw.Playlist, rawTracks))
		}
		ssplaylists = append(ssplaylists, plss)
	}
	return ssplaylists
}
