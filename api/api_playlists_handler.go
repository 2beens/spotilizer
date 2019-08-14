package api

import (
	"fmt"
	"io"
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
			handler.getPlaylistsSnapshots(user.Username, false, w)
		case "/api/ssplaylists/full":
			handler.getPlaylistsSnapshots(user.Username, true, w)
		default:
			handler.getPlaylistsSnapshot(user.Username, w, r)
		}
	case "DELETE":
		handler.deletePlaylistsSnapshot(user.Username, w, r)
	default:
		util.SendAPIErrorResp(w, "unknown/unsupported request method", http.StatusBadRequest)
	}
}

func (handler *PlaylistsHandler) getPlaylistsSnapshot(username string, w io.Writer, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Debugf(" > get playlists snapshot [%s]: username [%s]", timestamp, username)

	snapshotRaw, err := services.UserPlaylist.GetPlaylistsSnapshotByTimestamp(username, timestamp)
	if err != nil {
		log.Errorf(" >>> error while trying to get playlists snapshot: %s", err.Error())
		util.SendAPIErrorResp(w, "Error occurred: "+err.Error(), http.StatusNotFound)
		return
	}
	if snapshotRaw == nil {
		log.Errorf(" >>> error while trying to get playlists snapshot: snapshot is nil")
		util.SendAPIErrorResp(w, "Playlists snapshot not found", http.StatusNotFound)
		return
	}

	snapshots := handler.preparePlaylistsSnapshots([]models.PlaylistsSnapshot{*snapshotRaw}, true)
	if len(snapshots) == 0 {
		log.Errorf(" >>> error while trying to get playlists snapshot: DTO transformation error")
		util.SendAPIErrorResp(w, "internal server error", http.StatusInternalServerError)
		return
	}

	util.SendAPIOKRespWithData(w, "success", snapshots[0])
}

func (handler *PlaylistsHandler) deletePlaylistsSnapshot(username string, w io.Writer, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Debugf(" > delete playlists snapshot [%s]: username [%s]", timestamp, username)

	snapshot, err := services.UserPlaylist.DeletePlaylistsSnapshot(username, timestamp)
	if err != nil {
		log.Errorf(" >>> error while trying to delete playlists snapshot: %s", err.Error())
		util.SendAPIErrorResp(w, "Error occurred: "+err.Error(), http.StatusNotFound)
		return
	}
	if snapshot == nil {
		log.Errorf(" >>> error while trying to delete playlists snapshot: snapshot is nil")
		util.SendAPIErrorResp(w, "Playlists snapshot not deleted: not found ", http.StatusNotFound)
		return
	}

	util.SendAPIOKResp(w, fmt.Sprintf("Playlists snapshot [%s] successfully deleted.", snapshot.Timestamp))
}

func (handler *PlaylistsHandler) getPlaylistsSnapshots(username string, loadAllData bool, w io.Writer) {
	log.Debugf(" > get playlists snapshots: username [%s]", username)
	ssplaylistsRaw := services.UserPlaylist.GetAllPlaylistsSnapshots(username)
	ssplaylists := handler.preparePlaylistsSnapshots(ssplaylistsRaw, loadAllData)
	util.SendAPIOKRespWithData(w, "success", ssplaylists)
}

func (handler *PlaylistsHandler) preparePlaylistsSnapshots(ssplaylistsRaw []models.PlaylistsSnapshot, loadTracks bool) []models.DTOPlaylistSnapshot {
	var ssplaylists []models.DTOPlaylistSnapshot
	for _, plssRaw := range ssplaylistsRaw {
		plss := models.DTOPlaylistSnapshot{
			Timestamp: plssRaw.Timestamp.Unix(),
			Playlists: []models.DTOPlaylist{},
		}
		var rawTracks []models.SpPlaylistTrack
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
