package services

import (
	"encoding/json"
	"fmt"
	"log"

	c "github.com/2beens/spotilizer/config"
	m "github.com/2beens/spotilizer/models"
)

type UserPlaylistService interface {
	GetCurrentUserPlaylists(authOptions *m.SpotifyAuthOptions) (response m.SpGetCurrentPlaylistsResp, err error)
	GetSavedTracks(authOptions *m.SpotifyAuthOptions) (tracks []m.SpAddedTrack, err error)
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

// GetCurrentUserPlaylists more info: https://developer.spotify.com/console/get-current-user-playlists/
func (ups *SpotifyUserPlaylistService) GetCurrentUserPlaylists(authOptions *m.SpotifyAuthOptions) (playlists []m.SpPlaylist, err *m.SpAPIError) {
	offset := 0
	prevCount := 0
	for {
		path := fmt.Sprintf("%s?offset=%d&limit=50", ups.urlCurrentUserPlaylists, offset)
		body, err := getFromSpotify(ups.spotifyApiURL, path, authOptions)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error getting current user playlists. details: %v", err)
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}

		if apiErr, isError := getAPIError(body); isError {
			log.Printf(" >>> API error: status [%d] -> [%s]\n", apiErr.Error.Status, apiErr.Error.Message)
			return nil, &apiErr
		}

		var response m.SpGetCurrentPlaylistsResp
		err = json.Unmarshal(body, &response)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error occured while unmarshaling get playlists response: %v", err)
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}

		playlists = append(playlists, response.Items...)

		if len(response.Next) == 0 {
			return playlists, nil
		}

		// safety mechanism agains infinite loop - if no new tracks are added, bail out
		if prevCount == len(playlists) {
			log.Println(" > no new tracks coming in, bail out")
			return playlists, nil
		}
		prevCount = len(playlists)

		offset += 50
	}
}

func (ups *SpotifyUserPlaylistService) GetSavedTracks(authOptions *m.SpotifyAuthOptions) (tracks []m.SpAddedTrack, err *m.SpAPIError) {
	offset := 0
	prevCount := 0
	for {
		path := fmt.Sprintf("%s?offset=%d&limit=50", ups.urlCurrentUserSavedTracks, offset)
		body, err := getFromSpotify(ups.spotifyApiURL, path, authOptions)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error getting current user tracks. details: %v", err)
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}

		if apiErr, isError := getAPIError(body); isError {
			log.Printf(" >>> API error: status [%d] -> [%s]\n", apiErr.Error.Status, apiErr.Error.Message)
			return nil, &apiErr
		}

		var response m.SpGetSavedTracksResp
		err = json.Unmarshal(body, &response)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error occured while unmarshaling get tracks response: %v", err)
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}

		tracks = append(tracks, response.Items...)

		if len(response.Next) == 0 {
			return tracks, nil
		}

		// safety mechanism agains infinite loop - if no new tracks are added, bail out
		if prevCount == len(tracks) {
			log.Println(" > no new tracks coming in, bail out")
			return tracks, nil
		}
		prevCount = len(tracks)

		offset += 50
	}
}
