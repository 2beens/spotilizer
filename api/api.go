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

func DeletePlaylistsSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Debugln(" > deleting playlists: " + timestamp)
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		log.Errorf(" >>> user/cookie error while trying to delete playlists snapshot: %s", err.Error())
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Debugf(" > delete playlists snapshot [%s]: username [%s]", timestamp, user.Username)

	snapshot, err := services.UserPlaylist.DeletePlaylistsSnapshot(user.Username, timestamp)
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

func DeleteFavTracksSnapshots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Debugln(" > deleting fav tracks: " + timestamp)
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		log.Errorf(" >>> user/cookie error while trying to delete fav. tracks snapshot: %s", err.Error())
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Debugf(" > delete fav tracks snapshot [%s]: username [%s]", timestamp, user.Username)

	snapshot, err := services.UserPlaylist.DeleteFavTracksSnapshot(user.Username, timestamp)
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

func GetPlaylistsSnapshotsHandler(loadAllData bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugln(" > API: getting user playists snapshots ...")
		user, err := services.Users.GetUserByRequestCookieID(r)
		if err != nil {
			log.Errorf(" >>> user/cookie error while getting playlists snapshots: %s", err.Error())
			util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
			return
		}

		log.Debugf(" > get playlists snapshots: username [%s]", user.Username)
		ssplaylists := preparePlaylistsSnapshots(user.Username, loadAllData)
		util.SendAPIOKRespWithData(w, "success", ssplaylists)
	}
}

func preparePlaylistsSnapshots(username string, loadTracks bool) []models.DTOPlaylistSnapshot {
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

func GetFavTracksSnapshotsHandler(loadAllData bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugln(" > API: getting user fav tracks snapshots ...")
		user, err := services.Users.GetUserByRequestCookieID(r)
		if err != nil {
			log.Errorf(" >>> user/cookie error while getting fav tracks snapshots: %s", err.Error())
			util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
			return
		}

		log.WithFields(log.Fields{
			"loadAllData": loadAllData,
		}).Debugf(" > get fav tracks snapshots: username [%s]", user.Username)

		sstracksRaw := services.UserPlaylist.GetAllFavTracksSnapshots(user.Username)
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
}
