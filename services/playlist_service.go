package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	c "github.com/2beens/spotilizer/config"
	m "github.com/2beens/spotilizer/models"
)

type UserPlaylistService interface {
	GetCurrentUserPlaylists(authOptions m.SpotifyAuthOptions) (response m.SpGetCurrentPlaylistsResp, err error)
	GetSavedTracks(authOptions m.SpotifyAuthOptions) (tracks []m.SpAddedTrack, err error)
}

// TODO: removed this, it is unnecessary, especially that all these values can be found in config obj
type SpotifyUserPlaylistService struct {
	spotifyApiURL             string
	urlCurrentUserPlaylists   string
	urlCurrentUserSavedTracks string
}

func NewSpotifyUserPlaylistService() *SpotifyUserPlaylistService {
	var ps SpotifyUserPlaylistService
	ps.spotifyApiURL = c.Get().SpotifyApiURL
	ps.urlCurrentUserPlaylists = c.Get().URLCurrentUserPlaylists
	ps.urlCurrentUserSavedTracks = c.Get().URLCurrentUserSavedTracks
	return &ps
}

func getAPIError(body []byte) (spErr m.SpError, isError bool) {
	err := json.Unmarshal(body, &spErr)
	if err != nil {
		return m.SpError{}, false
	}
	return spErr, true
}

// GetCurrentUserPlaylists more info: https://developer.spotify.com/console/get-current-user-playlists/
func (ups SpotifyUserPlaylistService) GetCurrentUserPlaylists(authOptions m.SpotifyAuthOptions) (response m.SpGetCurrentPlaylistsResp, err error) {
	body, err := ups.getFromSpotify(ups.urlCurrentUserPlaylists, authOptions)
	if err != nil {
		log.Printf(" >>> error getting current user playlists. details: %v\n", err)
		return m.SpGetCurrentPlaylistsResp{}, err
	}
	json.Unmarshal(body, &response)
	return
}

func (ups SpotifyUserPlaylistService) GetSavedTracks(authOptions m.SpotifyAuthOptions) (tracks []m.SpAddedTrack, err error) {
	offset := 0
	prevCount := 0
	for {
		path := fmt.Sprintf("%s?offset=%d&limit=50", ups.urlCurrentUserSavedTracks, offset)
		body, err := ups.getFromSpotify(path, authOptions)
		if err != nil {
			log.Printf(" >>> error getting current user tracks. details: %v\n", err)
			return nil, err
		}
		var response m.SpGetSavedTracksResp
		err = json.Unmarshal(body, &response)
		if err != nil {
			if apiErr, isError := getAPIError(body); isError {
				return nil, fmt.Errorf("API error: [%s] -> [%s]\n", apiErr.Error, apiErr.ErrorDescription)
			}
			log.Printf(" >>> error occured while unmarshaling get tracks response: %v\n", err)
			return nil, err
		}

		tracks = append(tracks, response.Items...)

		if len(response.Next) == 0 {
			return tracks, nil
		}

		// safety mechanism agains infinite loop - if no new tracks are added, bail out
		if prevCount == len(tracks) {
			return tracks, nil
		}
		prevCount = len(tracks)

		offset += 50
	}
}

func (ups SpotifyUserPlaylistService) getFromSpotify(path string, authOptions m.SpotifyAuthOptions) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ups.spotifyApiURL+path, nil)
	if err != nil {
		log.Printf(" >>> error getting spotify response. details: %v\n", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authOptions.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf(" >>> error getting current user playlist. details: %v\n", err)
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(" >>> error getting spotify response. details: %v\n", err)
		return nil, err
	}

	// log.Println(string(body))

	return
}
