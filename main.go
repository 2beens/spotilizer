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
	h "github.com/2beens/spotilizer/handlers"
	m "github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/util"
	"github.com/gorilla/mux"
)

var serverURL = fmt.Sprintf("%s://%s:%s", c.Protocol, c.IPAddress, c.Port)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/index" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	util.RenderView(w, "index", m.ViewData{})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	util.RenderView(w, "contact", m.ViewData{})
}

// middleware function wrapping a handler functiomn and logging the request path
func logMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(" > request path: [%s]\n", r.URL.Path)
		f(w, r)
	}
}

func routerSetup() (r *mux.Router) {
	r = mux.NewRouter()

	// server static files
	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	// web content
	r.HandleFunc("/", logMiddleware(indexHandler))
	r.HandleFunc("/contact", logMiddleware(contactHandler))

	// router example usage with params (remove later)
	r.HandleFunc("/books/{title}/page/{page}", logMiddleware(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"] // the book title slug
		page := vars["page"]   // the page
		log.Printf(" > received title [%s] and page [%s]\n", title, page)
	})).Methods("GET")

	// spotify API
	r.HandleFunc("/login", logMiddleware(h.GetSpotifyLoginHandler(serverURL)))
	r.HandleFunc("/callback", logMiddleware(h.GetSpotifyCallbackHandler(serverURL)))
	r.HandleFunc("/refresh_token", logMiddleware(h.GetRefreshTokenHandler(serverURL)))
	r.HandleFunc("/save_current_playlists", logMiddleware(h.GetSaveCurrentPlaylistsHandler(serverURL)))
	r.HandleFunc("/save_current_tracks", logMiddleware(h.GetSaveCurrentTracksHandler(serverURL)))

	return
}

/****************** M A I N ************************************************************************/
/***************************************************************************************************/
func main() {
	displayHelp := flag.Bool("h", false, "display info/help message")
	logFileName := flag.String("logfile", "", "log file used to store server logs")
	flag.Parse()

	if *displayHelp {
		fmt.Println("\t -h \t\t\t\t> show this message\n\t -logfile=<logFileName> \t> output log file name")
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
	c := make(chan os.Signal, 1)
	// we'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	signal.Notify(c, os.Interrupt)

	// block until (eg. Ctrl+C) signal is received
	<-c

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
