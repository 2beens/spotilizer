package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/2beens/spotilizer/db"
	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

func GetPlaylistsSnapshots(w http.ResponseWriter, r *http.Request) {
	user, err := services.Users.GetUserByRequestCookieID(r)
	if err != nil {
		log.Printf(" >>> %s\n", fmt.Sprintf(" >>> user/cookie error while saving current user tracks: %s", err.Error()))
		util.SendAPIErrorResp(w, "Not available when logged off", http.StatusForbidden)
		return
	}

	log.Printf(" > get playlists snapshots: username [%s]\n", user.Username)

	ssplaylistsRaw := db.GetAllPlaylistsSnapshots(user.Username)
	ssplaylists := []models.DTOPlaylistSnapshot{}
	for _, plssRaw := range *ssplaylistsRaw {
		plss := models.DTOPlaylistSnapshot{
			Timestamp: plssRaw.Timestamp.Unix(),
			Playlists: []models.DTOPlaylist{},
		}
		for _, plRaw := range plssRaw.Playlists {
			pl := models.DTOPlaylist{
				ID:         plRaw.ID,
				Name:       plRaw.Name,
				URI:        plRaw.URI,
				TracksHref: plRaw.Tracks.Href,
				Tracks:     []models.DTOTrack{},
			}
			plss.Playlists = append(plss.Playlists, pl)
		}
		ssplaylists = append(ssplaylists, plss)
	}

	util.SendAPIOKRespWithData(w, "success", ssplaylists)
}
