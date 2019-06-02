package db

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	m "github.com/2beens/spotilizer/models"
	"gopkg.in/redis.v3"
)

type SpotifyDBClient interface {
	SaveFavTracksSnapshot(ft *m.FavTracksSnapshot) (saved bool)
	SavePlaylistsSnapshot(ps *m.PlaylistsSnapshot) (saved bool)
	GetFavTrakcsSnapshot(key string) *m.FavTracksSnapshot
	GetPlaylistSnapshot(key string) *m.PlaylistsSnapshot
	GetAllFavTracksSnapshots(username string) *[]m.FavTracksSnapshot
	GetAllPlaylistsSnapshots(username string) *[]m.PlaylistsSnapshot
}

type SpotifyDB struct{}

func (self SpotifyDB) SaveFavTracksSnapshot(ft *m.FavTracksSnapshot) (saved bool) {
	log.Printf(" > saving fav tracks [%d] for user [%s] ...\n", len(ft.Tracks), ft.Username)
	tracksJSON, err := json.Marshal(ft.Tracks)
	if err != nil {
		log.Println(" >>> json marshalling error saving tracks to DB for user: " + ft.Username)
		return false
	}
	timestamp := strconv.FormatInt(ft.Timestamp.Unix(), 10)
	snapshotKey := fmt.Sprintf("favtracksshot::user::%s::timestamp::%s", ft.Username, timestamp)
	log.Printf(" > saving new playlist snapshot: [%s]\n", snapshotKey)
	cmd := rc.Set(snapshotKey, string(tracksJSON), 0)
	if err := cmd.Err(); err != nil {
		log.Printf(" >>> failed to store tracks snapshot for user: %s\n", ft.Username)
		return false
	}
	log.Printf(" > user [%s] fav tracks snapshot saved to DB\n", ft.Username)
	return true
}

func (self SpotifyDB) SavePlaylistsSnapshot(ps *m.PlaylistsSnapshot) (saved bool) {
	log.Printf(" > saving playlists [%d] for user [%s] ...\n", len(ps.Playlists), ps.Username)
	playlistsJSON, err := json.Marshal(ps.Playlists)
	if err != nil {
		log.Println(" >>> json marshalling error saving playlists to DB for user: " + ps.Username)
		return false
	}
	timestamp := strconv.FormatInt(ps.Timestamp.Unix(), 10)
	snapshotKey := fmt.Sprintf("playlistsshot::user::%s::timestamp::%s", ps.Username, timestamp)
	log.Printf(" > saving new playlist snapshot: [%s]\n", snapshotKey)
	cmd := rc.Set(snapshotKey, string(playlistsJSON), 0)
	if err := cmd.Err(); err != nil {
		log.Printf(" >>> failed to store playlists snapshot for user: %s\n", ps.Username)
		return false
	}
	log.Printf(" > user [%s] playlists snapshot saved to DB\n", ps.Username)
	return true
}

func (self SpotifyDB) GetFavTrakcsSnapshot(key string) *m.FavTracksSnapshot {
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
		log.Println(" >>> error while parsing fav. tracks snapshot timestamp")
		return nil
	}
	timestamp := time.Unix(timestampInt, 0)

	tracksJSON := cmd.Val()
	tracks := &[]m.SpAddedTrack{}
	err = json.Unmarshal([]byte(tracksJSON), tracks)
	if err != nil {
		log.Printf(" >>> failed to unmarshal fav. tracks for snapshot [%s]: %s\n", key, err.Error())
		return nil
	}

	return &m.FavTracksSnapshot{Username: username, Timestamp: timestamp, Tracks: *tracks}
}

func (self SpotifyDB) GetPlaylistSnapshot(key string) *m.PlaylistsSnapshot {
	cmd := rc.Get(key)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get playlist snapshot [%s]: %s\n", key, err.Error())
		return nil
	}

	keyParts := strings.Split(key, "::")
	// log.Printf(" > get playlists snapshot, key parts: %v\n", keyParts)
	username := keyParts[2]
	timestampStr := keyParts[4]
	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		log.Println(" >>> error while parsing playlist snapshot timestamp")
		return nil
	}
	timestamp := time.Unix(timestampInt, 0)

	playlistsJSON := cmd.Val()
	playlists := &[]m.SpPlaylist{}
	err = json.Unmarshal([]byte(playlistsJSON), playlists)
	if err != nil {
		log.Printf(" >>> failed to unmarshal playlists for snapshot [%s]: %s\n", key, err.Error())
		return nil
	}

	return &m.PlaylistsSnapshot{Username: username, Timestamp: timestamp, Playlists: *playlists}
}

func (self SpotifyDB) GetAllFavTracksSnapshots(username string) *[]m.FavTracksSnapshot {
	snapshotsKey := fmt.Sprintf("favtracksshot::user::%s::timestamp::*", username)
	cmd := rc.Keys(snapshotsKey)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get all playlists snapshots for user [%s]: %s\n", username, err.Error())
		return nil
	}
	favtsnapshots := []m.FavTracksSnapshot{}
	for _, skey := range cmd.Val() {
		ft := self.GetFavTrakcsSnapshot(skey)
		if ft != nil {
			favtsnapshots = append(favtsnapshots, *ft)
		}
	}
	return &favtsnapshots
}

func (self SpotifyDB) GetAllPlaylistsSnapshots(username string) *[]m.PlaylistsSnapshot {
	snapshotsKey := fmt.Sprintf("playlistsshot::user::%s::timestamp::*", username)
	cmd := rc.Keys(snapshotsKey)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Printf(" >>> failed to get all playlists snapshots for user [%s]: %s\n", username, err.Error())
		return nil
	}
	plsnapshots := []m.PlaylistsSnapshot{}
	for _, skey := range cmd.Val() {
		ps := self.GetPlaylistSnapshot(skey)
		if ps != nil {
			plsnapshots = append(plsnapshots, *ps)
		}
	}
	return &plsnapshots
}
