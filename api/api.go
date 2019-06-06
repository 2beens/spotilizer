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
			// TODO: download tracks (maybe get a different API, not to send all data at once)
			// plTracks = append(plTracks, models.DTOTrack{URI: plRaw.Tracks.Href})
			// tracksBody, err := services.GetFromSpotify(plRaw.Tracks.Href, "", user.Auth)
			// if err != nil {
			// 	log.Printf(" >>> error, cannot get playlist tracks, user [%s], for playlist [%s]\n", user.Username, plRaw.Name)
			// 	continue
			// }
			// apiErr, isError := services.GetAPIError(tracksBody)
			// if isError {
			// 	// TODO: refresh token in case of expired (status 401, The access token expired)
			// 	log.Printf(" >>> error, cannot get playlist tracks, user [%s], for playlist [%s]. Error: [status %v] %s\n",
			// 		user.Username, plRaw.Name, apiErr.Error.Status, apiErr.Error.Message)
			// 	continue
			// }

			// playlistTracksRaw := &models.SpGetPlaylistTracksResp{}
			// err = json.Unmarshal(tracksBody, &playlistTracksRaw)
			// if err != nil {
			// 	log.Printf(" >>> error, cannot get playlist tracks, user [%s], for playlist [%s]\n", user.Username, plRaw.Name)
			// }

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
			Tracks:    []models.DTOAddedTrack{},
		}
		for _, trRaw := range tracksssRaw.Tracks {
			tracksss.Tracks = append(tracksss.Tracks, models.SpAddedTrack2dtoAddedTrack(trRaw))
		}
		sstracks = append(sstracks, tracksss)
	}

	util.SendAPIOKRespWithData(w, "success", sstracks)
}
