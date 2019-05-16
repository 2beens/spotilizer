package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var stateKey = "spotify_auth_state"
var port = "8080"

// spotify things
var clientID string
var clientSecret string
var loginRedirectURL = fmt.Sprintf("http://localhost:%s/callback", port)

// var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(" > request path: [%s]\n", r.URL.Path)
	if r.URL.Path != "/index" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	render(w, "index", ViewData{})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(" > request path: [%s]\n", r.URL.Path)
	render(w, "contact", ViewData{})
}

// templates cheatsheet
// https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet
func render(w http.ResponseWriter, page string, viewData ViewData) {
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
	state := generateRandomString(16)
	addCookie(w, stateKey, state)

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
		http.Redirect(w, r, "http://localhost:8080", 302)
		return
	}

	code, codeOk := q["code"]
	state, stateOk := q["state"]
	storedStateCookie, sStateCookieErr := r.Cookie(stateKey)
	if !codeOk || !stateOk {
		log.Println(" > login failed, error: some of the mandatory params not found")
		// TODO: redirect to some error, or show error on the index page
		http.Redirect(w, r, "http://localhost:8080", 302)
		return
	}

	if storedStateCookie == nil || storedStateCookie.Value != state[0] || sStateCookieErr != nil {
		log.Printf(" > login failed, error: state cookie not found or state mismatch. more details [%v]\n", err)
		if storedStateCookie != nil {
			log.Printf(" >>> storedStateCookie: [%s] state: [%s]\n", storedStateCookie.Value, state[0])
		} else {
			log.Println(" >>> storedStateCookie is nil!")
		}
		// TODO: redirect to some error, or show error on the index page
		http.Redirect(w, r, "http://localhost:8080", 302)
		return
	}

	cleearCookie(w, stateKey)
	authOptions := makeAuthPostReq(code[0])
	at := authOptions.AccessToken
	rt := authOptions.RefreshToken
	log.Printf(" > success! AT [%s] RT [%s]\n", at, rt)
	log.Printf(" > %v\n", authOptions)

	// redirect to index page with acces and refresh tokens
	render(w, "index", ViewData{Message: "success", Error: "", Data: authOptions})
}

// https://developer.spotify.com/documentation/general/guides/authorization-guide/
func makeAuthPostReq(code string) SpotifyAuthOptions {
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
		return SpotifyAuthOptions{}
	}
	defer resp.Body.Close()
	log.Println("------------------------------------------")
	log.Println("response Status:", resp.Status)
	// redirect to error if status != 200
	body, _ := ioutil.ReadAll(resp.Body)
	authOptions := SpotifyAuthOptions{}
	json.Unmarshal(body, &authOptions)
	return authOptions
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

	loggingSetup(*logFileName)

	// read spotify client ID & Secret
	var err error
	clientID, clientSecret, err = readSpotifyAuthData()
	if err != nil {
		log.Println(err)
		return
	}

	// server static files
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// index
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/contact", contactHandler)

	// spotify API
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", spotifyCallbackHandler)
	http.HandleFunc("/refresh_token", refreshTokenHandler)

	log.Printf(" > server listening on port: [%s]\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
