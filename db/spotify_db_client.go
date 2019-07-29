package db

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/2beens/spotilizer/models"
	"gopkg.in/redis.v3"
)

type SpotifyDBClient interface {
	SaveFavTracksSnapshot(ft *models.FavTracksSnapshot) (saved bool)
	SavePlaylistsSnapshot(ps *models.PlaylistsSnapshot) (saved bool)
	DeletePlaylistsSnapshot(username string, timestamp string) (*models.PlaylistsSnapshot, error)
	DeleteFavTracksSnapshot(username string, timestamp string) (*models.FavTracksSnapshot, error)
	GetPlaylistsSnapshotByTimestamp(username string, timestamp string) (*models.PlaylistsSnapshot, error)
	GetFavTracksSnapshotByTimestamp(username string, timestamp string) (*models.FavTracksSnapshot, error)
	GetPlaylistsSnapshot(key string) *models.PlaylistsSnapshot
	GetFavTracksSnapshot(key string) *models.FavTracksSnapshot
	GetAllFavTracksSnapshots(username string) []models.FavTracksSnapshot
	GetAllPlaylistsSnapshots(username string) []models.PlaylistsSnapshot
}

type SpotifyDB struct{}

func (sDB SpotifyDB) SaveFavTracksSnapshot(ft *models.FavTracksSnapshot) (saved bool) {
	log.Tracef(" > saving fav tracks [%d] for user [%s] ...\n", len(ft.Tracks), ft.Username)
	tracksJSON, err := json.Marshal(ft.Tracks)
	if err != nil {
		log.Println(" >>> json marshalling error saving tracks to DB for user: " + ft.Username)
		return false
	}
	timestamp := strconv.FormatInt(ft.Timestamp.Unix(), 10)
	snapshotKey := fmt.Sprintf("favtracksshot::user::%s::timestamp::%s", ft.Username, timestamp)
	log.Tracef(" > saving new playlist snapshot: [%s]\n", snapshotKey)
	cmd := rc.Set(snapshotKey, string(tracksJSON), 0)
	if err := cmd.Err(); err != nil {
		log.Printf(" >>> failed to store tracks snapshot for user: %s\n", ft.Username)
		return false
	}
	log.Debugf(" > user [%s] fav tracks snapshot saved to DB\n", ft.Username)
	return true
}

func (sDB SpotifyDB) SavePlaylistsSnapshot(ps *models.PlaylistsSnapshot) (saved bool) {
	log.Tracef(" > saving playlists [%d] for user [%s] ...\n", len(ps.Playlists), ps.Username)
	playlistsJSON, err := json.Marshal(ps.Playlists)
	if err != nil {
		log.Println(" >>> json marshalling error saving playlists to DB for user: " + ps.Username)
		return false
	}
	timestamp := strconv.FormatInt(ps.Timestamp.Unix(), 10)
	snapshotKey := fmt.Sprintf("playlistsshot::user::%s::timestamp::%s", ps.Username, timestamp)
	log.Tracef(" > saving new playlist snapshot: [%s]\n", snapshotKey)
	cmd := rc.Set(snapshotKey, string(playlistsJSON), 0)
	if err := cmd.Err(); err != nil {
		log.Printf(" >>> failed to store playlists snapshot for user: %s\n", ps.Username)
		return false
	}
	log.Debugf(" > user [%s] playlists snapshot saved to DB\n", ps.Username)
	return true
}

func (sDB SpotifyDB) DeletePlaylistsSnapshot(username string, timestamp string) (*models.PlaylistsSnapshot, error) {
	log.Tracef(" > deleting playlist snapshot [%s] ...\n", timestamp)
	snapshotKey := fmt.Sprintf("playlistsshot::user::%s::timestamp::%s", username, timestamp)
	snapshot := sDB.GetPlaylistsSnapshot(snapshotKey)
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot [%s] not found", timestamp)
	}

	cmd := rc.Del(snapshotKey)
	if err := cmd.Err(); err != nil {
		log.Debugf(" >>> failed to delete playlists snapshot [%s] for user [%s]: %s\n", timestamp, username, err.Error())
		return nil, err
	}

	deletedRecordsCount := cmd.Val()
	if deletedRecordsCount == 0 {
		return snapshot, fmt.Errorf("snapshot [%s] found, but not deleted", snapshot.Timestamp)
	}

	return snapshot, nil
}

func (sDB SpotifyDB) DeleteFavTracksSnapshot(username string, timestamp string) (*models.FavTracksSnapshot, error) {
	log.Tracef(" > deleting fav tracks snapshot [%s] ...\n", timestamp)
	snapshotKey := fmt.Sprintf("favtracksshot::user::%s::timestamp::%s", username, timestamp)
	snapshot := sDB.GetFavTracksSnapshot(snapshotKey)
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot [%s] not found", timestamp)
	}

	cmd := rc.Del(snapshotKey)
	if err := cmd.Err(); err != nil {
		log.Debugf(" >>> failed to delete fav tracks snapshot [%s] for user [%s]: %s\n", timestamp, username, err.Error())
		return nil, err
	}

	deletedRecordsCount := cmd.Val()
	if deletedRecordsCount == 0 {
		return snapshot, fmt.Errorf("snapshot [%s] found, but not deleted", snapshot.Timestamp)
	}

	return snapshot, nil
}

func (sDB SpotifyDB) GetPlaylistsSnapshotByTimestamp(username string, timestamp string) (*models.PlaylistsSnapshot, error) {
	log.Tracef(" > getting playlists snapshot [%s] ...\n", timestamp)
	snapshotKey := fmt.Sprintf("playlistsshot::user::%s::timestamp::%s", username, timestamp)
	snapshot := sDB.GetPlaylistsSnapshot(snapshotKey)
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot [%s] not found", timestamp)
	}
	return snapshot, nil
}

func (sDB SpotifyDB) GetFavTracksSnapshotByTimestamp(username string, timestamp string) (*models.FavTracksSnapshot, error) {
	log.Tracef(" > getting fav tracks snapshot [%s] ...\n", timestamp)
	snapshotKey := fmt.Sprintf("favtracksshot::user::%s::timestamp::%s", username, timestamp)
	snapshot := sDB.GetFavTracksSnapshot(snapshotKey)
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot [%s] not found", timestamp)
	}
	return snapshot, nil
}

func (sDB SpotifyDB) GetFavTracksSnapshot(key string) *models.FavTracksSnapshot {
	cmd := rc.Get(key)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get fav tracks snapshot [%s]: %s\n", key, err.Error())
		return nil
	}

	keyParts := strings.Split(key, "::")
	// log.Printf(" > get tracks snapshot, key parts: %v\n", keyParts)
	username := keyParts[2]
	timestampStr := keyParts[4]
	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		log.Debugf(" >>> error while parsing fav. tracks snapshot timestamp")
		return nil
	}
	timestamp := time.Unix(timestampInt, 0)

	tracksJSON := cmd.Val()
	tracks := &[]models.SpAddedTrack{}
	err = json.Unmarshal([]byte(tracksJSON), tracks)
	if err != nil {
		log.Errorf(" >>> failed to unmarshal fav. tracks for snapshot [%s]: %s\n", key, err.Error())
		return nil
	}

	return &models.FavTracksSnapshot{Username: username, Timestamp: timestamp, Tracks: *tracks}
}

func (sDB SpotifyDB) GetPlaylistsSnapshot(key string) *models.PlaylistsSnapshot {
	cmd := rc.Get(key)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Debugf(" >>> failed to get playlist snapshot [%s]: %s\n", key, err.Error())
		return nil
	}

	keyParts := strings.Split(key, "::")
	// log.Printf(" > get playlists snapshot, key parts: %v\n", keyParts)
	username := keyParts[2]
	timestampStr := keyParts[4]
	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		log.Debugf(" >>> error while parsing playlist snapshot timestamp")
		return nil
	}
	timestamp := time.Unix(timestampInt, 0)

	playlistsJSON := cmd.Val()
	playlists := &[]models.PlaylistSnapshot{}
	err = json.Unmarshal([]byte(playlistsJSON), playlists)
	if err != nil {
		log.Errorf(" >>> failed to unmarshal playlists for snapshot [%s]: %s\n", key, err.Error())
		return nil
	}

	return &models.PlaylistsSnapshot{Username: username, Timestamp: timestamp, Playlists: *playlists}
}

func (sDB SpotifyDB) GetAllFavTracksSnapshots(username string) []models.FavTracksSnapshot {
	snapshotsKey := fmt.Sprintf("favtracksshot::user::%s::timestamp::*", username)
	cmd := rc.Keys(snapshotsKey)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get all playlists snapshots for user [%s]: %s\n", username, err.Error())
		return nil
	}
	var favtsnapshots []models.FavTracksSnapshot
	for _, skey := range cmd.Val() {
		ft := sDB.GetFavTracksSnapshot(skey)
		if ft != nil {
			favtsnapshots = append(favtsnapshots, *ft)
		}
	}
	return favtsnapshots
}

func (sDB SpotifyDB) GetAllPlaylistsSnapshots(username string) []models.PlaylistsSnapshot {
	snapshotsKey := fmt.Sprintf("playlistsshot::user::%s::timestamp::*", username)
	cmd := rc.Keys(snapshotsKey)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get all playlists snapshots for user [%s]: %s\n", username, err.Error())
		return nil
	}
	var plsnapshots []models.PlaylistsSnapshot
	for _, skey := range cmd.Val() {
		ps := sDB.GetPlaylistsSnapshot(skey)
		if ps != nil {
			plsnapshots = append(plsnapshots, *ps)
		}
	}
	return plsnapshots
}
