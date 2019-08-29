package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/2beens/spotilizer/api"
	"github.com/2beens/spotilizer/constants"
	"github.com/2beens/spotilizer/db"
	"github.com/2beens/spotilizer/handlers"
	"github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	_ "net/http/pprof"
)

var serverURL = fmt.Sprintf("%s://%s:%s", constants.Protocol, constants.IPAddress, constants.Port)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cookieIDval string
		cookieID, err := r.Cookie(constants.CookieUserIDKey)
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
	r.HandleFunc("/", handlers.IndexHandler)
	r.HandleFunc("/about", handlers.AboutHandler)
	r.HandleFunc("/contact", handlers.ContactHandler)

	// spotify API
	r.HandleFunc("/login", handlers.GetSpotifyLoginHandler(serverURL))
	r.HandleFunc("/logout", handlers.LogoutHandler)
	r.HandleFunc("/callback", handlers.GetSpotifyCallbackHandler(serverURL))
	r.HandleFunc("/refresh_token", handlers.RefreshTokenHandler)
	r.HandleFunc("/save_current_playlists", handlers.SaveCurrentPlaylistsHandler)
	r.HandleFunc("/save_current_tracks", handlers.SaveCurrentTracksHandler)

	apiFavTracksHandler := api.NewFavTracksHandler(services.Users, services.UserPlaylist)
	apiPlaylistsHandler := api.NewPlaylistsHandler()

	r.Handle("/api/ssplaylists", apiPlaylistsHandler)
	r.Handle("/api/ssplaylists/full", apiPlaylistsHandler)
	r.Handle("/api/ssplaylists/{timestamp}", apiPlaylistsHandler)
	r.Handle("/api/ssfavtracks", apiFavTracksHandler)
	r.Handle("/api/ssfavtracks/full", apiFavTracksHandler)
	r.Handle("/api/ssfavtracks/{timestamp}", apiFavTracksHandler)
	// diffs
	r.Handle("/api/ssplaylists/diff/{timestamp}", apiPlaylistsHandler)
	r.Handle("/api/ssfavtracks/diff/{timestamp}", apiFavTracksHandler)

	// debugging
	r.HandleFunc("/debug", handlers.DebugHandler)

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

	// read and set view templates from disk
	util.SetupTemplates()

	// we need a webserver to get the pprof webserver
	pprofhost := "localhost"
	pprofport := "6060"
	go func() {
		log.Debugf("starting pprof server on [%s:%s] ...", pprofhost, pprofport)
		log.Debugln(http.ListenAndServe(pprofhost+":"+pprofport, nil))
	}()

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
	handlers.SetClientIDAndSecret(clientID, clientSecret)

	// redis setup
	db.InitRedisClient(*flashDB)
	// services setup
	services.InitServices()

	router := routerSetup()

	ipAndPort := fmt.Sprintf("%s:%s", constants.IPAddress, constants.Port)
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

	// channel for OS interrupt signal
	interruptCh := make(chan struct{}, 1)

	// starting CLI
	go cli(interruptCh)

	waitForInterruptSignal(interruptCh)
	gracefulShutdown(httpServer)
}

func cli(osInterruptCh chan struct{}) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		inputTxt, _ := reader.ReadString('\n')
		inputTxt = strings.Replace(inputTxt, "\n", "", -1)

		if inputTxt == "exit" {
			osInterruptCh <- struct{}{}
			fmt.Println("\n => will terminate the server ...")
		}
	}
}

func waitForInterruptSignal(interruptCh chan struct{}) {
	c := make(chan os.Signal, 1)
	// we'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	signal.Notify(c, os.Interrupt)
	// block until (eg. Ctrl+C) signal is received or EXIT command received
	select {
	case <-c:
		log.Warn(" > OS Interrupt signal received")
	case <-interruptCh:
		log.Warn(" > EXIT signal received")
	}
}

func gracefulShutdown(httpServer *http.Server) {
	log.Debug(" > graceful shutdown initiated ...")

	// store users cookies data
	services.Users.StoreCookiesToDB()

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

	log.Warn(" > server shut down")
	os.Exit(0)
}
