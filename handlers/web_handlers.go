package handlers

import (
	"log"
	"net/http"

	c "github.com/2beens/spotilizer/constants"
	"github.com/2beens/spotilizer/db"
	m "github.com/2beens/spotilizer/models"
	s "github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

func GetIndexHandler(username string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		playlists := db.GetAllPlaylistsSnapshots(username)
		// TODO: tracks can (maybe have) to be transfeterd to web client via API, not template
		// tracks := db.GetAllFavTracksSnapshots(username)
		// purpously use anonymous type/struct here, for learning purposes... for now
		util.RenderView(w, "index", m.ViewData{Username: username, Data: struct {
			PlaylistsSnapshots *[]m.PlaylistsSnapshot `json:"ssplaylists"`
			// TracksSnapshots    *[]m.FavTracksSnapshot `json:"sstracks"`
		}{
			playlists,
			// tracks,
		},
		})
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	username, _ := util.GetUsernameByRequestCookieID(r)
	GetIndexHandler(username)(w, r)
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	username, _ := util.GetUsernameByRequestCookieID(r)
	util.RenderView(w, "contact", m.ViewData{Username: username})
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	username, _ := util.GetUsernameByRequestCookieID(r)
	util.RenderView(w, "about", m.ViewData{Username: username})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookieID, err := r.Cookie(c.CookieUserIDKey)
	if err == nil {
		s.Users.RemoveUserCookie(cookieID.Value)
		util.CleearCookie(&w, c.CookieStateKey)
	}
	IndexHandler(w, r)
}

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	cookieID, _ := r.Cookie(c.CookieUserIDKey)
	user, _ := s.Users.GetUserByCookieID(cookieID.Value)
	log.Println("--------------- USER      ---------------------------------")
	log.Println(user)
	playlists := db.GetAllPlaylistsSnapshots(user.Username)
	log.Println("--------------- PLAYLISTS ---------------------------------")
	for _, p := range *playlists {
		log.Printf(" ====>>> [%v]: count %d\n", p.Timestamp, len(p.Playlists))
	}
	log.Println("--------------- TRACKS    ---------------------------------")
	favtracks := db.GetAllFavTracksSnapshots(user.Username)
	for _, t := range *favtracks {
		log.Printf(" ====>>> [%v]: count %d\n", t.Timestamp, len(t.Tracks))
	}
	log.Println("-------------------------------------------------------------")
	http.Redirect(w, r, "/", 302)
}
