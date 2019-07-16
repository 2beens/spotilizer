package services

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	c "github.com/2beens/spotilizer/config"
	"github.com/2beens/spotilizer/db"
	m "github.com/2beens/spotilizer/models"
)

type UserPlaylistService interface {
	DownloadCurrentUserPlaylists(authOptions *m.SpotifyAuthOptions) (response m.SpGetCurrentPlaylistsResp, err error)
	DownloadSavedFavTracks(authOptions *m.SpotifyAuthOptions) (tracks []m.SpAddedTrack, err error)
	GetFavTrakcsSnapshotByTimestamp(username string, timestamp string) (*m.FavTracksSnapshot, error)
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
	ps.spotifyAPIURL = c.Conf.SpotifyAPIURL
	ps.urlCurrentUserPlaylists = c.Conf.URLCurrentUserPlaylists
	ps.urlCurrentUserSavedTracks = c.Conf.URLCurrentUserSavedTracks
	return &ps
}

// DownloadCurrentUserPlaylists more info: https://developer.spotify.com/console/get-current-user-playlists/
func (ups *SpotifyUserPlaylistService) DownloadCurrentUserPlaylists(accessToken string) (playlists []m.SpPlaylist, err *m.SpAPIError) {
	offset := 0
	prevCount := 0
	for {
		path := fmt.Sprintf("%s?offset=%d&limit=50", ups.urlCurrentUserPlaylists, offset)
		body, err := getFromSpotify(ups.spotifyAPIURL, path, accessToken)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error getting current user playlists. details: %s", err.Error())
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}
		if apiErr, isError := getAPIError(body); isError {
			log.Printf(" >>> API error: status [%d] -> [%s]\n", apiErr.Error.Status, apiErr.Error.Message)
			return nil, &apiErr
		}

		var response m.SpGetCurrentPlaylistsResp
		err = json.Unmarshal(body, &response)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error occured while unmarshaling get playlists response: %s", err.Error())
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}

		playlists = append(playlists, response.Items...)

		if len(response.Next) == 0 {
			return playlists, nil
		}

		// safety mechanism against infinite loop - if no new tracks are added, bail out
		if prevCount == len(playlists) {
			log.Println(" > no new tracks coming in, bail out")
			return playlists, nil
		}
		prevCount = len(playlists)

		offset += 50
	}
}

func (ups *SpotifyUserPlaylistService) DownloadPlaylistTracks(accessToken string, href string, total int) (tracks []m.SpPlaylistTrack, err *m.SpAPIError) {
	tracks = []m.SpPlaylistTrack{}
	prevCount := 0
	nextHref := href
	for {
		body, err := getFromSpotify(nextHref, "", accessToken)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error getting playlist tracks. details: %s", err.Error())
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}
		if apiErr, isError := getAPIError(body); isError {
			log.Printf(" >>> API getting playlist tracks error: status [%d] -> [%s]\n", apiErr.Error.Status, apiErr.Error.Message)
			return nil, &apiErr
		}

		var response m.SpGetPlaylistTracksResp
		err = json.Unmarshal(body, &response)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error occured while unmarshaling get playlist tracks response: %s", err.Error())
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}

		tracks = append(tracks, response.Items...)
		if len(response.Next) == 0 || len(tracks) >= total {
			return tracks, nil
		}
		nextHref = response.Next

		// safety mechanism against infinite loop - if no new tracks are added, bail out
		if prevCount == len(tracks) {
			log.Println(" > no new tracks coming in, bail out")
			return tracks, nil
		}
		prevCount = len(tracks)
	}
}

func (ups *SpotifyUserPlaylistService) DownloadSavedFavTracks(accessToken string) (tracks []m.SpAddedTrack, err *m.SpAPIError) {
	offset := 0
	prevCount := 0
	for {
		path := fmt.Sprintf("%s?offset=%d&limit=50", ups.urlCurrentUserSavedTracks, offset)
		body, err := getFromSpotify(ups.spotifyAPIURL, path, accessToken)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error getting current user tracks. details: %s", err.Error())
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}
		if apiErr, isError := getAPIError(body); isError {
			log.Printf(" >>> API error: status [%d] -> [%s]\n", apiErr.Error.Status, apiErr.Error.Message)
			return nil, &apiErr
		}

		var response m.SpGetSavedTracksResp
		err = json.Unmarshal(body, &response)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error occured while unmarshaling get tracks response: %s", err.Error())
			return nil, &m.SpAPIError{Error: m.SpError{Status: 500, Message: errMsg}}
		}

		tracks = append(tracks, response.Items...)

		if len(response.Next) == 0 {
			return tracks, nil
		}

		// safety mechanism against infinite loop - if no new tracks are added, bail out
		if prevCount == len(tracks) {
			log.Println(" > no new tracks coming in, bail out")
			return tracks, nil
		}
		prevCount = len(tracks)

		offset += 50
	}
}

func (ups *SpotifyUserPlaylistService) SaveFavTracksSnapshot(ft *m.FavTracksSnapshot) (saved bool) {
	return ups.spotifyDB.SaveFavTracksSnapshot(ft)
}

func (ups *SpotifyUserPlaylistService) SavePlaylistsSnapshot(ps *m.PlaylistsSnapshot) (saved bool) {
	return ups.spotifyDB.SavePlaylistsSnapshot(ps)
}

func (ups *SpotifyUserPlaylistService) GetFavTracksSnapshotByTimestamp(username string, timestamp string) (*m.FavTracksSnapshot, error) {
	return ups.spotifyDB.GetFavTracksSnapshotByTimestamp(username, timestamp)
}

func (ups *SpotifyUserPlaylistService) GetPlaylistsSnapsotByTimestamp(username string, timestamp string) (*m.PlaylistsSnapshot, error) {
	return ups.spotifyDB.GetPlaylistsSnapsotByTimestamp(username, timestamp)
}

func (ups *SpotifyUserPlaylistService) GetAllFavTracksSnapshots(username string) []m.FavTracksSnapshot {
	return ups.spotifyDB.GetAllFavTracksSnapshots(username)
}

func (ups *SpotifyUserPlaylistService) GetAllPlaylistsSnapshots(username string) []m.PlaylistsSnapshot {
	return ups.spotifyDB.GetAllPlaylistsSnapshots(username)
}

func (ups *SpotifyUserPlaylistService) DeletePlaylistsSnapshot(username string, timestamp string) (*m.PlaylistsSnapshot, error) {
	return ups.spotifyDB.DeletePlaylistsSnapshot(username, timestamp)
}

func (ups *SpotifyUserPlaylistService) DeleteFavTracksSnapshot(username string, timestamp string) (*m.FavTracksSnapshot, error) {
	return ups.spotifyDB.DeleteFavTracksSnapshot(username, timestamp)
}
