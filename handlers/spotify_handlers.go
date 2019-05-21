package handlers

import (
	// b64 "encoding/base64"
	// "fmt"
	"log"
	"net/http"
	// "net/url"
	// "strings"

	c "github.com/2beens/spotilizer/constants"
	// m "github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	// "github.com/2beens/spotilizer/util"
)

// var ClientID string
// var ClientSecret string

func GetSaveCurrentPlaylistsHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie(c.CookieUserIDKey)
		if err != nil {
			// TOOD: redirect to error
			log.Printf(" >>> error while saving current user playlists: %v\n", err)
			return
		}

		log.Printf(" > user ID: %s\n", userID.Value)

		if authOptions, found := services.Users.User2authOptionsMap[userID.Value]; found {
			playlists, err := services.UserPlaylist.GetCurrentUserPlaylists(authOptions)
			if err != nil {
				log.Printf(" >>> error while saving current user playlists: %v\n", err)
				return
			}
			log.Printf(" > playlists count: %d\n", len(playlists.Items))
			// TODO: return standardized resp message
			w.Write([]byte("saved!"))
			return
		}
		log.Printf(" >>> failed to find user, must login first\n")
		http.Redirect(w, r, serverURL, 302)
	}
}

// func GetSpotifyCallbackHandler(serverURL string) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf(" > request path: [%s]\n", r.URL.Path)

// 		q := r.URL.Query()
// 		err, ok := q["error"]
// 		if ok {
// 			log.Printf(" > login failed, error: [%v]\n", err)
// 			// TODO: redirect to some error, or show error on the index page
// 			http.Redirect(w, r, serverURL, 302)
// 			return
// 		}

// 		code, codeOk := q["code"]
// 		state, stateOk := q["state"]
// 		storedStateCookie, sStateCookieErr := r.Cookie(c.CookieStateKey)
// 		if !codeOk || !stateOk {
// 			log.Println(" > login failed, error: some of the mandatory params not found")
// 			// TODO: redirect to some error, or show error on the index page
// 			http.Redirect(w, r, serverURL, 302)
// 			return
// 		}

// 		if storedStateCookie == nil || storedStateCookie.Value != state[0] || sStateCookieErr != nil {
// 			log.Printf(" > login failed, error: state cookie not found or state mismatch. more details [%v]\n", err)
// 			// TODO: redirect to some error, or show error on the index page
// 			http.Redirect(w, r, serverURL, 302)
// 			return
// 		}

// 		util.CleearCookie(w, c.CookieStateKey)
// 		authOptions := makeAuthPostReq(code[0], serverURL)
// 		at := authOptions.AccessToken
// 		rt := authOptions.RefreshToken
// 		log.Printf(" > success! AT [%s] RT [%s]\n", at, rt)
// 		log.Printf(" > %v\n", authOptions)

// 		newUserID := util.GenerateRandomString(35)
// 		services.Users.User2authOptionsMap[newUserID] = authOptions
// 		util.AddCookie(w, c.CookieUserIDKey, newUserID)

// 		// redirect to index page with acces and refresh tokens
// 		render(w, "index", m.ViewData{Message: "success", Error: "", Data: authOptions})
// 	}
// }

// // https://developer.spotify.com/documentation/general/guides/authorization-guide/
// func makeAuthPostReq(code string, serverURL string) m.SpotifyAuthOptions {
// 	apiURL := "https://accounts.spotify.com"
// 	resource := "/api/token/"
// 	data := url.Values{}
// 	data.Set("code", code)
// 	data.Set("redirect_uri", fmt.Sprintf("%s/callback", serverURL))
// 	data.Set("grant_type", "authorization_code")

// 	u, _ := url.ParseRequestURI(apiURL)
// 	u.Path = resource
// 	urlStr := u.String()

// 	client := &http.Client{}
// 	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
// 	authEncoding := b64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
// 	r.Header.Add("Authorization", "Basic "+authEncoding)
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

// 	resp, err := client.Do(r)
// 	if err != nil {
// 		log.Printf(" >>> error making an auth post req: %v\n", err)
// 		return m.SpotifyAuthOptions{}
// 	}
// 	defer resp.Body.Close()
// 	log.Println("------------------------------------------")
// 	log.Println("response Status:", resp.Status)
// 	// redirect to error if status != 200
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	authOptions := m.SpotifyAuthOptions{}
// 	json.Unmarshal(body, &authOptions)
// 	return authOptions
// }
