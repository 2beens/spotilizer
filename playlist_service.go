package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// https://developer.spotify.com/console/get-current-user-playlists/
const spotifyApiURL = "https://api.spotify.com"
const urlCurrentUserPlaylists = "/v1/me/playlists"

func getCurrentUserPlaylists(authOptions SpotifyAuthOptions) (response SpGetCurrentPlaylistsResp, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", spotifyApiURL+urlCurrentUserPlaylists, nil)
	if err != nil {
		log.Printf(" >>> error getting current user playlist. details: %v\n", err)
		return SpGetCurrentPlaylistsResp{}, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authOptions.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf(" >>> error getting current user playlist. details: %v\n", err)
		return SpGetCurrentPlaylistsResp{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(" >>> error getting current user playlist. details: %v\n", err)
		return SpGetCurrentPlaylistsResp{}, err
	}

	// log.Println(string(body))

	json.Unmarshal(body, &response)

	return
}
