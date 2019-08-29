package services

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/2beens/spotilizer/config"
	"github.com/2beens/spotilizer/db"
	"github.com/2beens/spotilizer/models"
)

type UserPlaylistService interface {
	DownloadCurrentUserPlaylists(accessToken string) (playlists []models.SpPlaylist, err *models.SpAPIError)
	DownloadPlaylistTracks(accessToken string, href string, total int) (tracks []models.SpPlaylistTrack, err *models.SpAPIError)
	DownloadSavedFavTracks(accessToken string) (tracks []models.SpAddedTrack, err *models.SpAPIError)
	SaveFavTracksSnapshot(ft *models.FavTracksSnapshot) (saved bool)
	SavePlaylistsSnapshot(ps *models.PlaylistsSnapshot) (saved bool)
	GetFavTracksSnapshotByTimestamp(username string, timestamp string) (*models.FavTracksSnapshot, error)
	GetPlaylistsSnapshotByTimestamp(username string, timestamp string) (*models.PlaylistsSnapshot, error)
	GetAllFavTracksSnapshots(username string) []models.FavTracksSnapshot
	GetAllPlaylistsSnapshots(username string) []models.PlaylistsSnapshot
	DeletePlaylistsSnapshot(username string, timestamp string) (*models.PlaylistsSnapshot, error)
	DeleteFavTracksSnapshot(username string, timestamp string) (*models.FavTracksSnapshot, error)
}

// TODO: removed this, it is unnecessary, especially that all these values can be found in config obj
type SpotifyUserPlaylistService struct {
	spotifyDB                 db.SpotifyDBClient
	spotifyAPIURL             string
	urlCurrentUserPlaylists   string
	urlCurrentUserSavedTracks string
}

func NewSpotifyUserPlaylistService(spotifyDB db.SpotifyDBClient) UserPlaylistService {
	ps := new(SpotifyUserPlaylistService)
	ps.spotifyDB = spotifyDB
	ps.spotifyAPIURL = config.Conf.SpotifyAPIURL
	ps.urlCurrentUserPlaylists = config.Conf.URLCurrentUserPlaylists
	ps.urlCurrentUserSavedTracks = config.Conf.URLCurrentUserSavedTracks
	return ps
}

// DownloadCurrentUserPlaylists more info: https://developer.spotify.com/console/get-current-user-playlists/
func (ups *SpotifyUserPlaylistService) DownloadCurrentUserPlaylists(accessToken string) (playlists []models.SpPlaylist, err *models.SpAPIError) {
	offset := 0
	prevCount := 0
	for {
		path := fmt.Sprintf("%s?offset=%d&limit=50", ups.urlCurrentUserPlaylists, offset)
		body, err := getFromSpotify(ups.spotifyAPIURL, path, accessToken)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error getting current user playlists. details: %s", err.Error())
			return nil, &models.SpAPIError{Error: models.SpError{Status: 500, Message: errMsg}}
		}
		if apiErr, isError := getAPIError(body); isError {
			log.Printf(" >>> API error: status [%d] -> [%s]\n", apiErr.Error.Status, apiErr.Error.Message)
			return nil, &apiErr
		}

		var response models.SpGetCurrentPlaylistsResp
		err = json.Unmarshal(body, &response)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error occurred while unmarshaling get playlists response: %s", err.Error())
			return nil, &models.SpAPIError{Error: models.SpError{Status: 500, Message: errMsg}}
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

func (ups *SpotifyUserPlaylistService) DownloadPlaylistTracks(accessToken string, href string, total int) (tracks []models.SpPlaylistTrack, err *models.SpAPIError) {
	tracks = []models.SpPlaylistTrack{}
	prevCount := 0
	nextHref := href
	for {
		body, err := getFromSpotify(nextHref, "", accessToken)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error getting playlist tracks. details: %s", err.Error())
			return nil, &models.SpAPIError{Error: models.SpError{Status: 500, Message: errMsg}}
		}
		if apiErr, isError := getAPIError(body); isError {
			log.Printf(" >>> API getting playlist tracks error: status [%d] -> [%s]\n", apiErr.Error.Status, apiErr.Error.Message)
			return nil, &apiErr
		}

		var response models.SpGetPlaylistTracksResp
		err = json.Unmarshal(body, &response)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error occurred while unmarshaling get playlist tracks response: %s", err.Error())
			return nil, &models.SpAPIError{Error: models.SpError{Status: 500, Message: errMsg}}
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

func (ups *SpotifyUserPlaylistService) DownloadSavedFavTracks(accessToken string) (tracks []models.SpAddedTrack, err *models.SpAPIError) {
	offset := 0
	prevCount := 0
	for {
		path := fmt.Sprintf("%s?offset=%d&limit=50", ups.urlCurrentUserSavedTracks, offset)
		body, err := getFromSpotify(ups.spotifyAPIURL, path, accessToken)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error getting current user tracks. details: %s", err.Error())
			return nil, &models.SpAPIError{Error: models.SpError{Status: 500, Message: errMsg}}
		}
		if apiErr, isError := getAPIError(body); isError {
			log.Printf(" >>> API error: status [%d] -> [%s]\n", apiErr.Error.Status, apiErr.Error.Message)
			return nil, &apiErr
		}

		var response models.SpGetSavedTracksResp
		err = json.Unmarshal(body, &response)
		if err != nil {
			errMsg := fmt.Sprintf(" >>> error occurred while unmarshaling get tracks response: %s", err.Error())
			return nil, &models.SpAPIError{Error: models.SpError{Status: 500, Message: errMsg}}
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

func (ups *SpotifyUserPlaylistService) SaveFavTracksSnapshot(ft *models.FavTracksSnapshot) (saved bool) {
	return ups.spotifyDB.SaveFavTracksSnapshot(ft)
}

func (ups *SpotifyUserPlaylistService) SavePlaylistsSnapshot(ps *models.PlaylistsSnapshot) (saved bool) {
	return ups.spotifyDB.SavePlaylistsSnapshot(ps)
}

func (ups *SpotifyUserPlaylistService) GetFavTracksSnapshotByTimestamp(username string, timestamp string) (*models.FavTracksSnapshot, error) {
	return ups.spotifyDB.GetFavTracksSnapshotByTimestamp(username, timestamp)
}

func (ups *SpotifyUserPlaylistService) GetPlaylistsSnapshotByTimestamp(username string, timestamp string) (*models.PlaylistsSnapshot, error) {
	return ups.spotifyDB.GetPlaylistsSnapshotByTimestamp(username, timestamp)
}

func (ups *SpotifyUserPlaylistService) GetAllFavTracksSnapshots(username string) []models.FavTracksSnapshot {
	return ups.spotifyDB.GetAllFavTracksSnapshots(username)
}

func (ups *SpotifyUserPlaylistService) GetAllPlaylistsSnapshots(username string) []models.PlaylistsSnapshot {
	return ups.spotifyDB.GetAllPlaylistsSnapshots(username)
}

func (ups *SpotifyUserPlaylistService) DeletePlaylistsSnapshot(username string, timestamp string) (*models.PlaylistsSnapshot, error) {
	return ups.spotifyDB.DeletePlaylistsSnapshot(username, timestamp)
}

func (ups *SpotifyUserPlaylistService) DeleteFavTracksSnapshot(username string, timestamp string) (*models.FavTracksSnapshot, error) {
	return ups.spotifyDB.DeleteFavTracksSnapshot(username, timestamp)
}
