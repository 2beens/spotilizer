package services

import (
	"encoding/json"
	"fmt"
	"log"

	c "github.com/2beens/spotilizer/config"
	"github.com/2beens/spotilizer/db"
	m "github.com/2beens/spotilizer/models"
)

type UserPlaylistService interface {
	DownloadCurrentUserPlaylists(authOptions *m.SpotifyAuthOptions) (response m.SpGetCurrentPlaylistsResp, err error)
	DownloadSavedFavTracks(authOptions *m.SpotifyAuthOptions) (tracks []m.SpAddedTrack, err error)
	GetAllFavTracksSnapshots(username string) *[]m.FavTracksSnapshot
	GetAllPlaylistsSnapshots(username string) *[]m.PlaylistsSnapshot
	SaveFavTracksSnapshot(ft *m.FavTracksSnapshot) (saved bool)
	SavePlaylistsSnapshot(ps *m.PlaylistsSnapshot) (saved bool)
}

// TODO: removed this, it is unnecessary, especially that all these values can be found in config obj
type SpotifyUserPlaylistService struct {
	spotifyDB                 db.SpotifyDBClient
	spotifyAPIURL             string
	urlCurrentUserPlaylists   string
	urlCurrentUserSavedTracks string
}

func NewSpotifyUserPlaylistService(spotifyDB db.SpotifyDBClient) *SpotifyUserPlaylistService {
	var ps SpotifyUserPlaylistService
	ps.spotifyDB = spotifyDB
	ps.spotifyAPIURL = c.Get().SpotifyApiURL
	ps.urlCurrentUserPlaylists = c.Get().URLCurrentUserPlaylists
	ps.urlCurrentUserSavedTracks = c.Get().URLCurrentUserSavedTracks
	return &ps
}

// DownloadCurrentUserPlaylists more info: https://developer.spotify.com/console/get-current-user-playlists/
func (ups SpotifyUserPlaylistService) DownloadCurrentUserPlaylists(authOptions *m.SpotifyAuthOptions) (playlists []m.SpPlaylist, err *m.SpAPIError) {
	offset := 0
	prevCount := 0
	for {
		path := fmt.Sprintf("%s?offset=%d&limit=50", ups.urlCurrentUserPlaylists, offset)
		body, err := getFromSpotify(ups.spotifyAPIURL, path, authOptions)
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

func (ups SpotifyUserPlaylistService) DownloadSavedFavTracks(authOptions *m.SpotifyAuthOptions) (tracks []m.SpAddedTrack, err *m.SpAPIError) {
	offset := 0
	prevCount := 0
	for {
		path := fmt.Sprintf("%s?offset=%d&limit=50", ups.urlCurrentUserSavedTracks, offset)
		body, err := getFromSpotify(ups.spotifyAPIURL, path, authOptions)
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

func (self SpotifyUserPlaylistService) SaveFavTracksSnapshot(ft *m.FavTracksSnapshot) (saved bool) {
	return self.spotifyDB.SaveFavTracksSnapshot(ft)
}

func (self SpotifyUserPlaylistService) SavePlaylistsSnapshot(ps *m.PlaylistsSnapshot) (saved bool) {
	return self.spotifyDB.SavePlaylistsSnapshot(ps)
}

func (self SpotifyUserPlaylistService) GetAllFavTracksSnapshots(username string) *[]m.FavTracksSnapshot {
	return self.spotifyDB.GetAllFavTracksSnapshots(username)
}

func (self SpotifyUserPlaylistService) GetAllPlaylistsSnapshots(username string) *[]m.PlaylistsSnapshot {
	return self.spotifyDB.GetAllPlaylistsSnapshots(username)
}
