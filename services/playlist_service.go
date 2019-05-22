package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	m "github.com/2beens/spotilizer/models"
)

type UserPlaylistService interface {
	GetCurrentUserPlaylists(authOptions m.SpotifyAuthOptions) (response m.SpGetCurrentPlaylistsResp, err error)
}

type SpotifyUserPlaylistService struct {
	spotifyApiURL           string
	urlCurrentUserPlaylists string
}

// GetCurrentUserPlaylists more info: https://developer.spotify.com/console/get-current-user-playlists/
func (ups SpotifyUserPlaylistService) GetCurrentUserPlaylists(authOptions m.SpotifyAuthOptions) (response m.SpGetCurrentPlaylistsResp, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ups.spotifyApiURL+ups.urlCurrentUserPlaylists, nil)
	if err != nil {
		log.Printf(" >>> error getting current user playlist. details: %v\n", err)
		return m.SpGetCurrentPlaylistsResp{}, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authOptions.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf(" >>> error getting current user playlist. details: %v\n", err)
		return m.SpGetCurrentPlaylistsResp{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(" >>> error getting current user playlist. details: %v\n", err)
		return m.SpGetCurrentPlaylistsResp{}, err
	}

	// log.Println(string(body))

	json.Unmarshal(body, &response)

	return
}
