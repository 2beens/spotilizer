package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
	"github.com/gorilla/mux"
)

func DeletePlaylistsSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Println(" > deleting playlists: " + timestamp)
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		log.Printf(" >>> user/cookie error while trying to delete playlists snapshot: %s\n", err.Error())
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Printf(" > delete playlists snapshot [%s]: username [%s]\n", timestamp, user.Username)

	snapshot, err := services.UserPlaylist.DeletePlaylistsSnapshot(user.Username, timestamp)
	if err != nil {
		log.Printf(" >>> error while trying to delete playlists snapshot: %s\n", err.Error())
		util.SendAPIErrorResp(w, "Error occured: "+err.Error(), http.StatusNotFound)
		return
	}
	if snapshot == nil {
		log.Println(" >>> error while trying to delete playlists snapshot: snapshot is nil")
		util.SendAPIErrorResp(w, "Playlists snapshot not deleted: not found ", http.StatusNotFound)
		return
	}

	util.SendAPIOKResp(w, fmt.Sprintf("Playlists snapshot [%s] successfully deleted.", snapshot.Timestamp))
}

func DeleteFavTracksSnapshots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timestamp := vars["timestamp"]
	log.Println(" > deleting fav tracks: " + timestamp)
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		log.Printf(" >>> user/cookie error while trying to delete fav. tracks snapshot: %s\n", err.Error())
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Printf(" > delete fav tracks snapshot [%s]: username [%s]\n", timestamp, user.Username)

	snapshot, err := services.UserPlaylist.DeleteFavTracksSnapshot(user.Username, timestamp)
	if err != nil {
		log.Printf(" >>> error while trying to delete fav. tracks snapshot: %s\n", err.Error())
		util.SendAPIErrorResp(w, "Error occured: "+err.Error(), http.StatusNotFound)
		return
	}
	if snapshot == nil {
		log.Println(" >>> error while trying to delete fav. tracks snapshot: snapshot is nil")
		util.SendAPIErrorResp(w, "Favorite tracks snapshot not deleted: not found ", http.StatusNotFound)
		return
	}

	util.SendAPIOKResp(w, fmt.Sprintf("Favorite tracks snapshot [%s] successfully deleted.", snapshot.Timestamp))
}

func GetPlaylistsSnapshots(w http.ResponseWriter, r *http.Request) {
	log.Println(" > API: getting user playists snapshots ...")
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		log.Printf(" >>> user/cookie error while getting playlists snapshots: %s\n", err.Error())
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Printf(" > get playlists snapshots: username [%s]\n", user.Username)

	ssplaylistsRaw := services.UserPlaylist.GetAllPlaylistsSnapshots(user.Username)
	ssplaylists := []models.DTOPlaylistSnapshot{}
	for _, plssRaw := range ssplaylistsRaw {
		plss := models.DTOPlaylistSnapshot{
			Timestamp: plssRaw.Timestamp.Unix(),
			Playlists: []models.DTOPlaylist{},
		}
		for _, plRaw := range plssRaw.Playlists {
			plss.Playlists = append(plss.Playlists, models.SpPlaylist2dtoPlaylist(plRaw.Playlist, plRaw.Tracks))
		}
		ssplaylists = append(ssplaylists, plss)
	}

	util.SendAPIOKRespWithData(w, "success", ssplaylists)
}

func GetFavTracksSnapshots(w http.ResponseWriter, r *http.Request) {
	log.Println(" > API: getting user fav tracks snapshots ...")
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		log.Printf(" >>> user/cookie error while getting fav tracks snapshots: %s\n", err.Error())
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Printf(" > get fav tracks snapshots: username [%s]\n", user.Username)

	sstracksRaw := services.UserPlaylist.GetAllFavTracksSnapshots(user.Username)
	sstracks := []models.DTOFavTracksSnapshot{}
	for _, tracksssRaw := range sstracksRaw {
		tracksss := models.DTOFavTracksSnapshot{
			Timestamp: tracksssRaw.Timestamp.Unix(),
			Tracks:    []models.DTOTrack{},
		}
		for _, trRaw := range tracksssRaw.Tracks {
			tracksss.Tracks = append(tracksss.Tracks, models.SpAddedTrack2dtoTrack(trRaw))
		}
		sstracks = append(sstracks, tracksss)
	}

	util.SendAPIOKRespWithData(w, "success", sstracks)
}
