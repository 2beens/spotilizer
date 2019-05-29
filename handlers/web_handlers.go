package handlers

import (
	"net/http"

	c "github.com/2beens/spotilizer/constants"
	m "github.com/2beens/spotilizer/models"
	s "github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	username, _ := util.GetUsernameByRequestCookieID(r)
	util.RenderView(w, "index", m.ViewData{Username: username})
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
