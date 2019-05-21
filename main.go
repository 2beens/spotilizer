package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	c "github.com/2beens/spotilizer/constants"
	h "github.com/2beens/spotilizer/handlers"
	m "github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	"github.com/2beens/spotilizer/util"
	"github.com/gorilla/mux"
)

var serverURL = fmt.Sprintf("%s://%s:%s", c.Protocol, c.IPAddress, c.Port)

// spotify things
var clientID string
var clientSecret string
var loginRedirectURL = fmt.Sprintf("%s/callback", serverURL)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(" > request path: [%s]\n", r.URL.Path)
	if r.URL.Path != "/index" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	render(w, "index", m.ViewData{})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(" > request path: [%s]\n", r.URL.Path)
	render(w, "contact", m.ViewData{})
}

// templates cheatsheet
// https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet
func render(w http.ResponseWriter, page string, viewData m.ViewData) {
	// TODO: parse the template once and reuse it
	files := []string{
		"public/views/layouts/layout.html",
		"public/views/layouts/footer.html",
		"public/views/layouts/navbar.html",
		"public/views/" + page + ".html",
	}
	t, err := template.New("layout").ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "layout", viewData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(" > request path: [%s]\n", r.URL.Path)
	state := util.GenerateRandomString(16)
	util.AddCookie(w, c.CookieStateKey, state)

	q := url.Values{}
	q.Add("response_type", "code")
	q.Add("client_id", clientID)
	q.Add("scope", "user-read-private user-read-email user-library-read")
	q.Add("redirect_uri", loginRedirectURL)
	q.Add("state", state)

	redirectURL := "https://accounts.spotify.com/authorize?" + q.Encode()
	log.Println(" > /login, redirect to: " + redirectURL)
	http.Redirect(w, r, redirectURL, 302)
}

func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(" > request path: [%s]\n", r.URL.Path)

	q := r.URL.Query()
	refreshToken, refreshTokenOk := q["refresh_token"]
	if !refreshTokenOk {
		log.Println(" > refresh token failed, error: refresh_token param not found")
		// TODO: redirect to some error, or show error on the index page
		w.Write([]byte("missing refresh_token param"))
		return
	}
	log.Println(" > refresh token, value: " + refreshToken[0])

	//TODO: implement the rest ...

}

func spotifyCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(" > request path: [%s]\n", r.URL.Path)

	q := r.URL.Query()
	err, ok := q["error"]
	if ok {
		log.Printf(" > login failed, error: [%v]\n", err)
		// TODO: redirect to some error, or show error on the index page
		http.Redirect(w, r, serverURL, 302)
		return
	}

	code, codeOk := q["code"]
	state, stateOk := q["state"]
	storedStateCookie, sStateCookieErr := r.Cookie(c.CookieStateKey)
	if !codeOk || !stateOk {
		log.Println(" > login failed, error: some of the mandatory params not found")
		// TODO: redirect to some error, or show error on the index page
		http.Redirect(w, r, serverURL, 302)
		return
	}

	if storedStateCookie == nil || storedStateCookie.Value != state[0] || sStateCookieErr != nil {
		log.Printf(" > login failed, error: state cookie not found or state mismatch. more details [%v]\n", err)
		// TODO: redirect to some error, or show error on the index page
		http.Redirect(w, r, serverURL, 302)
		return
	}

	util.CleearCookie(w, c.CookieStateKey)
	authOptions := makeAuthPostReq(code[0])
	at := authOptions.AccessToken
	rt := authOptions.RefreshToken
	log.Printf(" > success! AT [%s] RT [%s]\n", at, rt)
	log.Printf(" > %v\n", authOptions)

	newUserID := util.GenerateRandomString(35)
	services.Users.User2authOptionsMap[newUserID] = authOptions
	util.AddCookie(w, c.CookieUserIDKey, newUserID)

	// redirect to index page with acces and refresh tokens
	render(w, "index", m.ViewData{Message: "success", Error: "", Data: authOptions})
}

// https://developer.spotify.com/documentation/general/guides/authorization-guide/
func makeAuthPostReq(code string) m.SpotifyAuthOptions {
	apiURL := "https://accounts.spotify.com"
	resource := "/api/token/"
	data := url.Values{}
	data.Set("code", code)
	data.Set("redirect_uri", loginRedirectURL)
	data.Set("grant_type", "authorization_code")

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	authEncoding := b64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	r.Header.Add("Authorization", "Basic "+authEncoding)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		log.Printf(" >>> error making an auth post req: %v\n", err)
		return m.SpotifyAuthOptions{}
	}
	defer resp.Body.Close()
	log.Println("------------------------------------------")
	log.Println("response Status:", resp.Status)
	// redirect to error if status != 200
	body, _ := ioutil.ReadAll(resp.Body)
	authOptions := m.SpotifyAuthOptions{}
	json.Unmarshal(body, &authOptions)
	return authOptions
}

func routerSetup() (r *mux.Router) {
	// https://github.com/gorilla/mux
	r = mux.NewRouter()

	// server static files
	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	// index
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/contact", contactHandler)

	// router example usage with params
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"] // the book title slug
		page := vars["page"]   // the page
		log.Printf(" > received title [%s] and page [%s]\n", title, page)
	}).Methods("GET")

	// spotify API
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/callback", spotifyCallbackHandler)
	r.HandleFunc("/refresh_token", refreshTokenHandler)
	r.HandleFunc("/save_current_playlists", h.GetSaveCurrentPlaylistsHandler(serverURL))

	return
}

// realy nice site on creating web applications in go:
// https://gowebexamples.com/routes-using-gorilla-mux/
// serving static files with go:
// https://www.alexedwards.net/blog/serving-static-sites-with-go
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
	var err error
	clientID, clientSecret, err = util.ReadSpotifyAuthData()
	if err != nil {
		log.Println(err)
		return
	}

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
