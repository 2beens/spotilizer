package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var port = "8080"

// templates ****
var templates = template.Must(template.ParseGlob("public/views/*.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func loadPage(title string) (*Page, error) {
	filename := "pages/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathMatch := validPath.FindStringSubmatch(r.URL.Path)
		if pathMatch == nil {
			fmt.Printf(" > path not found: %v\n", pathMatch)
			http.NotFound(w, r)
			return
		}
		fmt.Printf(" > path found: %v\n", pathMatch)
		fn(w, r, pathMatch[2])
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(" > request path: [%s]\n", r.URL.Path)
	if r.URL.Path != "/index" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	err := templates.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(" > request path: [%s]\n", r.URL.Path)
	err := templates.ExecuteTemplate(w, "contact.html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var stateKey = "spotify_auth_state"

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(" > request path: [%s]\n", r.URL.Path)
	state := generateRandomString(16)
	addCookie(w, stateKey, state)

	q := url.Values{}
	q.Add("response_type", "code")
	q.Add("client_id", clientID)
	q.Add("scope", "user-read-private user-read-email")
	q.Add("redirect_uri", loginRedirectURL)
	q.Add("state", state)

	redirectURL := "https://accounts.spotify.com/authorize?" + q.Encode()
	fmt.Println(" > /login, redirect to: " + redirectURL)
	http.Redirect(w, r, redirectURL, 302)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(" > request path: [%s]\n", r.URL.Path)

	q := r.URL.Query()
	err, ok := q["error"]
	if ok {
		fmt.Printf(" > login failed, error: [%v]\n", err)
		// TODO: redirect to some error, or show error on the index page
		http.Redirect(w, r, "http://localhost:8080", 302)
		return
	}

	code, codeOk := q["code"]
	state, stateOk := q["state"]
	storedStateCookie, sStateCookieErr := r.Cookie(stateKey)
	if !codeOk || !stateOk {
		fmt.Println(" > login failed, error: some of the mandatory params not found")
		// TODO: redirect to some error, or show error on the index page
		http.Redirect(w, r, "http://localhost:8080", 302)
		return
	}

	if storedStateCookie == nil || storedStateCookie.Value != state[0] || sStateCookieErr != nil {
		fmt.Printf(" > login failed, error: state cookie not found or state mismatch. more details [%v]\n", err)
		if storedStateCookie != nil {
			fmt.Printf(" >>> storedStateCookie: [%s] state: [%s]\n", storedStateCookie.Value, state[0])
		} else {
			fmt.Println(" >>> storedStateCookie is nil!")
		}
		// TODO: redirect to some error, or show error on the index page
		http.Redirect(w, r, "http://localhost:8080", 302)
		return
	}

	cleearCookie(w, stateKey)
	authOptions := makeAuthPostReq(code[0])
	at := authOptions.AccessToken
	rt := authOptions.RefreshToken
	fmt.Printf(" > success! AT [%s] RT [%s]\n", at, rt)
	fmt.Printf(" > %v\n", authOptions)

	// TODO: redirect to index page with acces and refresh tokens

	http.Redirect(w, r, "http://localhost:8080", 302)
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
		fmt.Printf(" >>> error making an auth post req: %v\n", err)
		return SpotifyAuthOptions{}
	}
	defer resp.Body.Close()
	fmt.Println("------------------------------------------")
	fmt.Println("response Status:", resp.Status)
	// redirect to error if status != 200
	body, _ := ioutil.ReadAll(resp.Body)
	authOptions := SpotifyAuthOptions{}
	json.Unmarshal(body, &authOptions)
	return authOptions
}

var clientID string
var clientSecret string
var loginRedirectURL = "http://localhost:8080/callback"

func readSpotifyAuthData() error {
	clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	fmt.Println(" > client ID: " + clientID)
	fmt.Println(" > client secret: " + clientSecret)
	if clientID == "" {
		return errors.New(" >>> error, client ID missing. set it using env [SPOTIFY_CLIENT_ID]")
	}
	if clientSecret == "" {
		return errors.New(" >>> error, client secret missing. set it using env [SPOTIFY_CLIENT_SECRET]")
	}
	return nil
}

// realy nice site on creating web applications in go:
// https://gowebexamples.com/routes-using-gorilla-mux/
// serving static files with go:
// https://www.alexedwards.net/blog/serving-static-sites-with-go
func main() {
	err := readSpotifyAuthData()
	if err != nil {
		fmt.Println(err)
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
	http.HandleFunc("/callback", callbackHandler)

	// will be removed later
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	fmt.Printf(" > server listening on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
