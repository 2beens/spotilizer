package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	c "github.com/2beens/spotilizer/constants"
	db "github.com/2beens/spotilizer/db"
	h "github.com/2beens/spotilizer/handlers"
	s "github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
	"github.com/gorilla/mux"
)

var serverURL = fmt.Sprintf("%s://%s:%s", c.Protocol, c.IPAddress, c.Port)

// middleware function wrapping a handler functiomn and logging the request path
func middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cookieIDval string
		cookieID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			cookieIDval = "<nil>"
		} else {
			cookieIDval = cookieID.Value
		}
		log.Printf(" ====> request path: [%s], cookieID: [%s]\n", r.URL.Path, cookieIDval)
		f(w, r)
	}
}

func routerSetup() (r *mux.Router) {
	r = mux.NewRouter()

	// server static files
	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	// web content
	r.HandleFunc("/", middleware(h.IndexHandler))
	r.HandleFunc("/about", middleware(h.AboutHandler))
	r.HandleFunc("/contact", middleware(h.ContactHandler))

	// router example usage with params (remove later)
	r.HandleFunc("/books/{title}/page/{page}", middleware(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"] // the book title slug
		page := vars["page"]   // the page
		log.Printf(" > received title [%s] and page [%s]\n", title, page)
	})).Methods("GET")

	// spotify API
	r.HandleFunc("/login", middleware(h.GetSpotifyLoginHandler(serverURL)))
	r.HandleFunc("/logout", middleware(h.LogoutHandler))
	r.HandleFunc("/callback", middleware(h.GetSpotifyCallbackHandler(serverURL)))
	r.HandleFunc("/refresh_token", middleware(h.GetRefreshTokenHandler(serverURL)))
	r.HandleFunc("/save_current_playlists", middleware(h.GetSaveCurrentPlaylistsHandler(serverURL)))
	r.HandleFunc("/save_current_tracks", middleware(h.GetSaveCurrentTracksHandler(serverURL)))

	return
}

/****************** M A I N ************************************************************************/
/***************************************************************************************************/
func main() {
	displayHelp := flag.Bool("h", false, "display info/help message")
	flashDB := flag.Bool("flushdb", false, "Flush Redis DB")
	logFileName := flag.String("logfile", "", "log file used to store server logs")
	flag.Parse()

	if *displayHelp {
		fmt.Println(`
			-h                      > show this message
			-logfile=<logFileName>  > output log file name
			-flushdb                > flush/clear redis DB before start`)
		fmt.Println()
		return
	}

	util.LoggingSetup(*logFileName)

	// read spotify client ID & Secret
	clientID, clientSecret, err := util.ReadSpotifyAuthData()
	if err != nil {
		log.Println(err)
		return
	}
	h.SetCliendIdAndSecret(clientID, clientSecret)

	// redis setup
	db.InitRedisClient(*flashDB)
	// services setup
	s.InitServices()

	router := routerSetup()

	ipAndPort := fmt.Sprintf("%s:%s", c.IPAddress, c.Port)
	srv := &http.Server{
		Handler:      router,
		Addr:         ipAndPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// run our server in a goroutine so that it doesn't block
	go func() {
		log.Printf(" > server listening on: [%s]\n", ipAndPort)
		log.Fatal(srv.ListenAndServe())
	}()

	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	// TODO: maybe persist everything to redis db before shutdown ?

	c := make(chan os.Signal, 1)
	// we'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	signal.Notify(c, os.Interrupt)

	// block until (eg. Ctrl+C) signal is received
	<-c

	// store users cookies data
	s.Users.StoreCookiesInfo()

	// the duration for which the server gracefully wait for existing connections to finish
	maxWaitDuration := time.Second * 15
	// create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), maxWaitDuration)
	defer cancel()
	// doesn't block if no connections, but will otherwise wait until the timeout deadline
	srv.Shutdown(ctx)

	log.Println(" > shutting down")
	os.Exit(0)
}
