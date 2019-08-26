package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/2beens/spotilizer/constants"
	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"

	"github.com/stretchr/testify/assert"
)

func TestGetFavTracks(t *testing.T) {
	favTracksHandler := getTestFavTracksHandler()

	t.Run("test getting fav tracks without all data", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/ssfavtracks", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.AddCookie(&http.Cookie{
			Name:  constants.CookieUserIDKey,
			Value: "cookietu1",
		})
		resp := httptest.NewRecorder()
		favTracksHandler.ServeHTTP(resp, req)

		assert.NotNil(t, resp.Body)
		fmt.Printf(" >> resp body: %v\n", resp.Body)
		apiResp := checkFavTracksAPIResponse(t, resp.Body.Bytes())
		assert.Equal(t, 1, len(apiResp.Data))
		favTracksSnapshot := apiResp.Data[0]
		assert.Equal(t, 2, favTracksSnapshot.TracksCount)
		// do not load data -> tracks details have to be empty
		assert.Equal(t, 0, len(favTracksSnapshot.Tracks))
	})

	t.Run("test getting fav tracks with all data", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/ssfavtracks/full", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.AddCookie(&http.Cookie{
			Name:  constants.CookieUserIDKey,
			Value: "cookietu1",
		})
		resp := httptest.NewRecorder()
		favTracksHandler.ServeHTTP(resp, req)

		assert.NotNil(t, resp.Body)
		fmt.Printf(" >> resp body: %v\n", resp.Body)
		apiResp := checkFavTracksAPIResponse(t, resp.Body.Bytes())
		assert.Equal(t, 1, len(apiResp.Data))
		favTracksSnapshot := apiResp.Data[0]
		assert.Equal(t, 2, favTracksSnapshot.TracksCount)
		// load all data -> tracks details have to be here
		assert.Equal(t, 2, len(favTracksSnapshot.Tracks))
	})
}

func checkFavTracksAPIResponse(t *testing.T, rawResp []byte) *FavTracksAPIResponse {
	apiResp := &FavTracksAPIResponse{}
	err := json.Unmarshal(rawResp, apiResp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, apiResp.Status)
	assert.Equal(t, "success", apiResp.Message)
	assert.NotNil(t, apiResp.Data, "API response data must not be nil")
	return apiResp
}

func getTestFavTracksHandler() *FavTracksHandler {
	testUser, favTracksSnapshots := getFavTracksSnapshotsTestData()

	userPlaylistSrv := services.NewUserPlaylistTestService(favTracksSnapshots)
	testUserSrv := services.NewUserServiceTest()
	testUserSrv.Add(testUser)
	testUserSrv.AddUserCookie("cookietu1", testUser.Username)

	return NewFavTracksHandler(testUserSrv, userPlaylistSrv)
}

func getFavTracksSnapshotsTestData() (*models.User, []models.FavTracksSnapshot) {
	testUser := &models.User{
		Username: "testUser1",
		Auth: &models.SpotifyAuthOptions{
			AccessToken:  "testat",
			RefreshToken: "testrt",
		},
	}

	var ft1tracks []models.SpAddedTrack
	tr1 := models.SpAddedTrack{
		AddedAt: time.Date(2019, time.August, 1, 11, 0, 0, 0, time.UTC),
		Track: models.SpTrack{
			ID: "ft1tr1",
			Artists: []models.SpArtist{
				{
					ID:   "ft1art1",
					Name: "favTrack1 Artist",
					Type: "ft1art1type",
					Href: "dummy href 1",
				},
			},
			Explicit: true,
			Type:     "ft1tr1type",
			Album: models.SpAlbum{
				ID: "ft1tr1al1",
			},
			Name:        "favTrack1",
			TrackNumber: 1,
		},
	}
	tr2 := models.SpAddedTrack{
		AddedAt: time.Date(2019, time.August, 1, 10, 0, 0, 0, time.UTC),
		Track: models.SpTrack{
			ID: "ft1tr2",
			Artists: []models.SpArtist{
				{
					ID:   "ft1art2",
					Name: "favTrack2 Artist",
					Type: "ft1art2type",
					Href: "dummy href 2",
				},
			},
			Explicit: false,
			Type:     "ft1tr2type",
			Album: models.SpAlbum{
				ID: "ft1tr2al1",
			},
			Name:        "favTrack2",
			TrackNumber: 5,
		},
	}
	ft1tracks = append(ft1tracks, tr1)
	ft1tracks = append(ft1tracks, tr2)

	var favTracksSnapshots []models.FavTracksSnapshot
	ft1 := models.FavTracksSnapshot{
		Username:  testUser.Username,
		Timestamp: time.Date(2019, time.August, 1, 12, 0, 0, 0, time.UTC),
		Tracks:    ft1tracks,
	}
	favTracksSnapshots = append(favTracksSnapshots, ft1)

	return testUser, favTracksSnapshots
}

type FavTracksAPIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		Timestamp   int `json:"timestamp"`
		TracksCount int `json:"tracks_count"`
		Tracks      []struct {
			AddedAt int    `json:"added_at"`
			AddedBy string `json:"added_by"`
			Artists []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				Href string `json:"href"`
			} `json:"artists"`
			Album       string `json:"album"`
			URI         string `json:"uri"`
			ID          string `json:"id"`
			TrackNumber int    `json:"track_number"`
			DurationMs  int    `json:"duration_ms"`
			Name        string `json:"name"`
		} `json:"tracks"`
	} `json:"data"`
}
