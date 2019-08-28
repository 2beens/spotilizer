package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/2beens/spotilizer/constants"
	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"
	"github.com/gorilla/mux"

	"github.com/stretchr/testify/suite"
)

type FavTracksTestSuite struct {
	suite.Suite
	testUser  *models.User
	handler   *FavTracksHandler
	cookie    *http.Cookie
	snapshots []models.FavTracksSnapshot
}

func (suite *FavTracksTestSuite) SetupSuite() {
	suite.testUser = &models.User{
		Username: "testUser1",
		Auth: &models.SpotifyAuthOptions{
			AccessToken:  "test_accTok",
			RefreshToken: "test_refTok",
		},
	}

	suite.cookie = &http.Cookie{
		Name:  constants.CookieUserIDKey,
		Value: "cookietu1",
	}

	suite.fillSnapshotsTestData()

	testUserSrv := services.NewUserServiceTest()
	testUserSrv.Add(suite.testUser)
	testUserSrv.AddUserCookie("cookietu1", suite.testUser.Username)
	userPlaylistSrv := services.NewUserPlaylistTestService(suite.snapshots)

	suite.handler = NewFavTracksHandler(testUserSrv, userPlaylistSrv)
}

func (suite *FavTracksTestSuite) TestGetFavTracksCounts() {
	testCasePaths := []string{"/api/ssfavtracks", "/api/ssfavtracks/full", "/api/ssfavtracks"}
	for _, path := range testCasePaths {
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			suite.T().Fatal(err)
		}
		req.AddCookie(suite.cookie)

		resp := httptest.NewRecorder()
		suite.handler.ServeHTTP(resp, req)

		suite.NotNil(resp.Body)
		apiResp := suite.checkFavTracksAPIResponse(resp.Body.Bytes())

		suite.Equal(2, len(apiResp.Data))
		suite.Equal(2, apiResp.Data[0].TracksCount)
		suite.Equal(1, apiResp.Data[1].TracksCount)
		if strings.HasSuffix(path, "/full") {
			suite.Equal(2, len(apiResp.Data[0].Tracks))
			suite.Equal(1, len(apiResp.Data[1].Tracks))
		} else {
			suite.Equal(0, len(apiResp.Data[0].Tracks))
			suite.Equal(0, len(apiResp.Data[1].Tracks))
		}
	}
}

func (suite *FavTracksTestSuite) TestGetFavTracksDetails() {
	req, err := http.NewRequest("GET", "/api/ssfavtracks/full", nil)
	if err != nil {
		suite.FailNowf("TestGetFavTracksDetails error", "details: %s", err.Error())
	}
	req.AddCookie(suite.cookie)

	resp := httptest.NewRecorder()
	suite.handler.ServeHTTP(resp, req)

	suite.NotNil(resp.Body)
	apiResp := suite.checkFavTracksAPIResponse(resp.Body.Bytes())
	suite.Equal(2, len(apiResp.Data))

	orgSnp1 := suite.snapshots[0]
	favTracksSnapshot1 := apiResp.Data[0]
	snapshot1timestamp := time.Unix(int64(favTracksSnapshot1.Timestamp), 0)
	suite.True(snapshot1timestamp.Equal(orgSnp1.Timestamp))
	suite.Equal(len(orgSnp1.Tracks), favTracksSnapshot1.TracksCount, "wrong tracks count in fav tracks snapshot 1")
	orgTr1 := orgSnp1.Tracks[0]
	suite.Equal(orgTr1.Track.ID, favTracksSnapshot1.Tracks[0].ID)
	suite.Equal(orgTr1.Track.Artists[0].Name, favTracksSnapshot1.Tracks[0].Artists[0].Name)
	suite.Equal(orgTr1.Track.Name, favTracksSnapshot1.Tracks[0].Name)

	orgSnp2 := suite.snapshots[1]
	favTracksSnapshot2 := apiResp.Data[1]
	snapshot2timestamp := time.Unix(int64(favTracksSnapshot2.Timestamp), 0)
	suite.True(snapshot2timestamp.Equal(orgSnp2.Timestamp))
	suite.Equal(len(orgSnp2.Tracks), favTracksSnapshot2.TracksCount, "wrong tracks count in fav tracks snapshot 2")
	orgTr2 := orgSnp2.Tracks[0]
	suite.Equal(orgTr2.Track.ID, favTracksSnapshot2.Tracks[0].ID)
	suite.Equal(orgTr2.Track.Artists[0].Name, favTracksSnapshot2.Tracks[0].Artists[0].Name)
	suite.Equal(orgTr2.Track.Name, favTracksSnapshot2.Tracks[0].Name)
}

func (suite *FavTracksTestSuite) TestGetFavTracksByTimestamp() {
	timestamp := suite.snapshots[0].Timestamp.Unix()
	req := suite.getRequest("/api/ssfavtracks/")
	// TODO: needs a fix
	req = mux.SetURLVars(req, map[string]string{"timespan": strconv.FormatInt(timestamp, 10)})
	req.AddCookie(suite.cookie)

	resp := httptest.NewRecorder()
	suite.handler.ServeHTTP(resp, req)

	suite.NotNil(resp.Body)
	apiResp := suite.checkFavTracksAPIResponse(resp.Body.Bytes())
	suite.NotNil(apiResp)

	// TODO: assert the rest of the data
}

// In order for 'go test' to run this suite, we need to create a normal test function and pass our suite to suite.Run
func TestFavTracksTestSuite(t *testing.T) {
	suite.Run(t, new(FavTracksTestSuite))
}

func (suite *FavTracksTestSuite) getRequest(path string) *http.Request {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		suite.T().Fatal(err)
	}
	return req
}

func (suite *FavTracksTestSuite) checkFavTracksAPIResponse(rawResp []byte) *FavTracksAPIResponse {
	apiResp := &FavTracksAPIResponse{}
	err := json.Unmarshal(rawResp, apiResp)
	if err != nil {
		suite.FailNowf("fail to unmarshal FavTracksAPIResponse", "Detals: %s", err.Error())
	}
	suite.Equal(200, apiResp.Status)
	suite.Equal("success", apiResp.Message)
	suite.NotNil(apiResp.Data, "API response data must not be nil")
	return apiResp
}

func (suite *FavTracksTestSuite) fillSnapshotsTestData() {
	var ft1tracks []models.SpAddedTrack
	var ft2tracks []models.SpAddedTrack
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
	tr3 := models.SpAddedTrack{
		AddedAt: time.Date(2019, time.July, 28, 10, 0, 0, 0, time.UTC),
		Track: models.SpTrack{
			ID: "ft2tr1",
			Artists: []models.SpArtist{
				{
					ID:   "ft2art1",
					Name: "favTrack3 Artist",
					Type: "ft1art1type",
					Href: "dummy href 3",
				},
			},
			Explicit: false,
			Type:     "ft2tr1type",
			Album: models.SpAlbum{
				ID: "ft2tr1al1",
			},
			Name:        "favTrack3",
			TrackNumber: 5,
		},
	}
	ft1tracks = append(ft1tracks, tr1)
	ft1tracks = append(ft1tracks, tr2)
	ft2tracks = append(ft2tracks, tr3)

	ft1 := models.FavTracksSnapshot{
		Username:  suite.testUser.Username,
		Timestamp: time.Date(2019, time.August, 1, 12, 0, 0, 0, time.UTC),
		Tracks:    ft1tracks,
	}
	ft2 := models.FavTracksSnapshot{
		Username:  suite.testUser.Username,
		Timestamp: time.Date(2019, time.August, 2, 12, 0, 0, 0, time.UTC),
		Tracks:    ft2tracks,
	}
	suite.snapshots = append(suite.snapshots, ft1)
	suite.snapshots = append(suite.snapshots, ft2)
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
