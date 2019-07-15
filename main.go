package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/2beens/spotilizer/api"
	c "github.com/2beens/spotilizer/constants"
	db "github.com/2beens/spotilizer/db"
	h "github.com/2beens/spotilizer/handlers"
	s "github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var serverURL = fmt.Sprintf("%s://%s:%s", c.Protocol, c.IPAddress, c.Port)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cookieIDval string
		cookieID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			cookieIDval = "<nil>"
		} else {
			cookieIDval = cookieID.Value
		}
		log.Printf(" ====> request path: [%s], cookieID: [%s]\n", r.URL.Path, cookieIDval)
		// call the next handler, which can be another middleware in the chain, or the final handler
		next.ServeHTTP(w, r)
	})
}

func routerSetup() (r *mux.Router) {
	r = mux.NewRouter()

	// server static files
	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	// web content
	r.HandleFunc("/", h.IndexHandler)
	r.HandleFunc("/about", h.AboutHandler)
	r.HandleFunc("/contact", h.ContactHandler)

	// spotify API
	r.HandleFunc("/login", h.GetSpotifyLoginHandler(serverURL))
	r.HandleFunc("/logout", h.LogoutHandler)
	r.HandleFunc("/callback", h.GetSpotifyCallbackHandler(serverURL))
	r.HandleFunc("/refresh_token", h.RefreshTokenHandler)
	r.HandleFunc("/save_current_playlists", h.SaveCurrentPlaylistsHandler)
	r.HandleFunc("/save_current_tracks", h.SaveCurrentTracksHandler)

	r.HandleFunc("/api/ssplaylists", api.GetPlaylistsSnapshotsHandler(false))
	r.HandleFunc("/api/ssplaylists/full", api.GetPlaylistsSnapshotsHandler(true))
	r.HandleFunc("/api/ssplaylists/{timestamp}", api.DeletePlaylistsSnapshot).Methods("DELETE")
	r.HandleFunc("/api/ssfavtracks", api.GetFavTracksSnapshotsHandler(false))
	r.HandleFunc("/api/ssfavtracks/full", api.GetFavTracksSnapshotsHandler(true))
	r.HandleFunc("/api/ssfavtracks/{timestamp}", api.GetFavTracksSnapshot).Methods("GET")
	r.HandleFunc("/api/ssfavtracks/{timestamp}", api.DeleteFavTracksSnapshots).Methods("DELETE")

	// debuging
	r.HandleFunc("/debug", h.DebugHandler)

	// middleware
	r.Use(loggingMiddleware)

	return r
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
	// logrus has seven logging levels:
	//		Trace, Debug, Info, Warning, Error, Fatal, Panic
	log.SetLevel(log.TraceLevel)

	// logging example with fields
	// log.WithFields(log.Fields{
	// 	"omg":    true,
	// 	"number": 122,
	// }).Warn("The group's number increased tremendously!")

	// read spotify client ID & Secret
	clientID, clientSecret, err := util.ReadSpotifyAuthData()
	if err != nil {
		log.Fatal(err)
	}
	h.SetCliendIDAndSecret(clientID, clientSecret)

	// redis setup
	db.InitRedisClient(*flashDB)
	// services setup
	s.InitServices()

	router := routerSetup()

	ipAndPort := fmt.Sprintf("%s:%s", c.IPAddress, c.Port)
	httpServer := &http.Server{
		Handler:      router,
		Addr:         ipAndPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// run our server in a goroutine so that it doesn't block
	go func() {
		log.Infof(" > server listening on: [%s]", ipAndPort)
		log.Fatal(httpServer.ListenAndServe())
	}()

	go func() {
		if logFileName != nil && len(*logFileName) > 0 {
			// log output is set to file already, bail out
			return
		}
		c := make(chan os.Signal, 1)
		// SIGHUP signal is sent when a program loses its controlling terminal
		signal.Notify(c, syscall.SIGHUP)
		<-c
		util.LoggingSetup("serverlog")
		log.Warn(" > controlling terminal lost, logging switched to file [serverlog.log]")
	}()

	gracefulShutdown(httpServer)
}

func gracefulShutdown(httpServer *http.Server) {
	c := make(chan os.Signal, 1)
	// we'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	signal.Notify(c, os.Interrupt)

	// block until (eg. Ctrl+C) signal is received
	<-c

	// store users cookies data
	s.Users.StoreCookiesToDB()

	// the duration for which the server gracefully wait for existing connections to finish
	maxWaitDuration := time.Second * 15
	// create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), maxWaitDuration)
	defer cancel()
	// doesn't block if no connections, but will otherwise wait until the timeout deadline
	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Error(" >>> failed to gracefully shutdown")
	}

	log.Info(" > server shut down")
	os.Exit(0)
}
