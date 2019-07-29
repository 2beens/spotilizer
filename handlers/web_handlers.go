package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/2beens/spotilizer/constants"
	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

func GetIndexHandler(username string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		util.RenderView(w, "index", models.ViewData{Username: username})
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	username, _ := util.GetUsernameByRequestCookieID(r)
	GetIndexHandler(username)(w, r)
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	username, _ := util.GetUsernameByRequestCookieID(r)
	util.RenderView(w, "contact", models.ViewData{Username: username})
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	username, _ := util.GetUsernameByRequestCookieID(r)
	util.RenderView(w, "about", models.ViewData{Username: username})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookieID, err := r.Cookie(constants.CookieUserIDKey)
	if err == nil {
		services.Users.RemoveUserCookie(cookieID.Value)
		util.ClearCookie(&w, constants.CookieStateKey)
	}
	IndexHandler(w, r)
}

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	cookieID, _ := r.Cookie(constants.CookieUserIDKey)
	user, _ := services.Users.GetUserByCookieID(cookieID.Value)
	log.Infoln("--------------- USER      ---------------------------------")
	log.Infoln(user)
	playlists := services.UserPlaylist.GetAllPlaylistsSnapshots(user.Username)
	log.Infoln("--------------- PLAYLISTS ---------------------------------")
	for _, p := range playlists {
		log.Infof(" ====>>> [%v]: count %d\n", p.Timestamp, len(p.Playlists))
	}
	log.Infoln("--------------- TRACKS    ---------------------------------")
	favtracks := services.UserPlaylist.GetAllFavTracksSnapshots(user.Username)
	for _, t := range favtracks {
		log.Infof(" ====>>> [%v]: count %d\n", t.Timestamp, len(t.Tracks))
	}
	log.Infoln("-------------------------------------------------------------")
	http.Redirect(w, r, "/", http.StatusFound)
}
