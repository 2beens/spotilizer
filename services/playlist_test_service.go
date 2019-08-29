package services

import (
	"strconv"

	"github.com/2beens/spotilizer/models"
)

type UserPlaylistTestService struct {
	tracksSnapshots []models.FavTracksSnapshot
}

func NewUserPlaylistTestService(tracksSnapshots []models.FavTracksSnapshot) UserPlaylistService {
	return &UserPlaylistTestService{
		tracksSnapshots: tracksSnapshots,
	}
}

func (ups *UserPlaylistTestService) DownloadCurrentUserPlaylists(accessToken string) (playlists []models.SpPlaylist, err *models.SpAPIError) {
	return nil, nil
}

func (ups *UserPlaylistTestService) DownloadPlaylistTracks(accessToken string, href string, total int) (tracks []models.SpPlaylistTrack, err *models.SpAPIError) {
	return nil, nil
}

func (ups *UserPlaylistTestService) DownloadSavedFavTracks(accessToken string) (tracks []models.SpAddedTrack, err *models.SpAPIError) {
	return nil, nil
}

func (ups *UserPlaylistTestService) SaveFavTracksSnapshot(ft *models.FavTracksSnapshot) (saved bool) {
	return true
}

func (ups *UserPlaylistTestService) SavePlaylistsSnapshot(ps *models.PlaylistsSnapshot) (saved bool) {
	return true
}

func (ups *UserPlaylistTestService) GetFavTracksSnapshotByTimestamp(username string, timestamp string) (*models.FavTracksSnapshot, error) {
	for _, s := range ups.tracksSnapshots {
		s := s
		if strconv.FormatInt(s.Timestamp.Unix(), 10) == timestamp {
			return &s, nil
		}
	}
	return nil, nil
}

func (ups *UserPlaylistTestService) GetPlaylistsSnapshotByTimestamp(username string, timestamp string) (*models.PlaylistsSnapshot, error) {
	return nil, nil
}

func (ups *UserPlaylistTestService) GetAllFavTracksSnapshots(username string) []models.FavTracksSnapshot {
	return ups.tracksSnapshots
}

func (ups *UserPlaylistTestService) GetAllPlaylistsSnapshots(username string) []models.PlaylistsSnapshot {
	return nil
}

func (ups *UserPlaylistTestService) DeletePlaylistsSnapshot(username string, timestamp string) (*models.PlaylistsSnapshot, error) {
	return nil, nil
}

func (ups *UserPlaylistTestService) DeleteFavTracksSnapshot(username string, timestamp string) (*models.FavTracksSnapshot, error) {
	return nil, nil
}
