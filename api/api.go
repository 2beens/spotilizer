package api

import (
	"log"
	"net/http"

	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

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
