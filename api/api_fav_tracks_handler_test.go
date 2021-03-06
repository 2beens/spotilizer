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
	allTracks []models.SpAddedTrack
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
	userPlaylistSrv := services.NewUserPlaylistTestService(suite.allTracks, suite.snapshots)

	suite.handler = NewFavTracksHandler(testUserSrv, userPlaylistSrv)
}

func (suite *FavTracksTestSuite) TestGetAllFavTracksSnapshotsCounts() {
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
		apiResp := suite.checkAllFavTracksSnapshotsAPIResponse(resp.Body.Bytes())

		suite.Equal(2, len(apiResp.Snapshots))
		suite.Equal(2, apiResp.Snapshots[0].TracksCount)
		suite.Equal(1, apiResp.Snapshots[1].TracksCount)
		if strings.HasSuffix(path, "/full") {
			suite.Equal(2, len(apiResp.Snapshots[0].Tracks))
			suite.Equal(1, len(apiResp.Snapshots[1].Tracks))
		} else {
			suite.Equal(0, len(apiResp.Snapshots[0].Tracks))
			suite.Equal(0, len(apiResp.Snapshots[1].Tracks))
		}
	}
}

func (suite *FavTracksTestSuite) TestGetAllFavTracksSnapshotsDetails() {
	req, err := http.NewRequest("GET", "/api/ssfavtracks/full", nil)
	if err != nil {
		suite.FailNowf("TestGetFavTracksDetails error", "details: %s", err.Error())
	}
	req.AddCookie(suite.cookie)

	resp := httptest.NewRecorder()
	suite.handler.ServeHTTP(resp, req)

	suite.NotNil(resp.Body)
	apiResp := suite.checkAllFavTracksSnapshotsAPIResponse(resp.Body.Bytes())
	suite.Equal(2, len(apiResp.Snapshots))

	orgSnp1 := suite.snapshots[0]
	favTracksSnapshot1 := apiResp.Snapshots[0]
	snapshot1timestamp := time.Unix(int64(favTracksSnapshot1.Timestamp), 0)
	suite.True(snapshot1timestamp.Equal(orgSnp1.Timestamp))
	suite.Equal(len(orgSnp1.Tracks), favTracksSnapshot1.TracksCount, "wrong tracks count in fav tracks snapshot 1")
	orgTr1 := orgSnp1.Tracks[0]
	suite.Equal(orgTr1.Track.ID, favTracksSnapshot1.Tracks[0].ID)
	suite.Equal(orgTr1.Track.Artists[0].Name, favTracksSnapshot1.Tracks[0].Artists[0].Name)
	suite.Equal(orgTr1.Track.Name, favTracksSnapshot1.Tracks[0].Name)

	orgSnp2 := suite.snapshots[1]
	favTracksSnapshot2 := apiResp.Snapshots[1]
	snapshot2timestamp := time.Unix(int64(favTracksSnapshot2.Timestamp), 0)
	suite.True(snapshot2timestamp.Equal(orgSnp2.Timestamp))
	suite.Equal(len(orgSnp2.Tracks), favTracksSnapshot2.TracksCount, "wrong tracks count in fav tracks snapshot 2")
	orgTr2 := orgSnp2.Tracks[0]
	suite.Equal(orgTr2.Track.ID, favTracksSnapshot2.Tracks[0].ID)
	suite.Equal(orgTr2.Track.Artists[0].Name, favTracksSnapshot2.Tracks[0].Artists[0].Name)
	suite.Equal(orgTr2.Track.Name, favTracksSnapshot2.Tracks[0].Name)
}

func (suite *FavTracksTestSuite) TestGetFavTracksSnapshotByTimestamp() {
	orgSnapshot := suite.snapshots[0]
	req := suite.getRequest("/api/ssfavtracks/{timestamp}")
	req = mux.SetURLVars(req, map[string]string{"timestamp": strconv.FormatInt(orgSnapshot.Timestamp.Unix(), 10)})
	req.AddCookie(suite.cookie)

	resp := httptest.NewRecorder()
	suite.handler.ServeHTTP(resp, req)

	suite.NotNil(resp.Body)
	apiResp := suite.checkFavTracksSnapshotAPIResponse(resp.Body.Bytes())
	suite.NotNil(apiResp)

	recSnapshot := apiResp.Snapshot
	suite.Equal(len(orgSnapshot.Tracks), recSnapshot.TracksCount)
	suite.Equal(orgSnapshot.Tracks[0].AddedAt.Unix(), int64(recSnapshot.Tracks[0].AddedAt))
	suite.Equal(orgSnapshot.Tracks[1].AddedAt.Unix(), int64(recSnapshot.Tracks[1].AddedAt))
	suite.Equal(orgSnapshot.Tracks[0].Track.ID, recSnapshot.Tracks[0].ID)
	suite.Equal(orgSnapshot.Tracks[1].Track.ID, recSnapshot.Tracks[1].ID)
	suite.Equal(orgSnapshot.Tracks[0].Track.Name, recSnapshot.Tracks[0].Name)
	suite.Equal(orgSnapshot.Tracks[1].Track.Name, recSnapshot.Tracks[1].Name)
	suite.Equal(orgSnapshot.Tracks[0].Track.Artists[0].Name, recSnapshot.Tracks[0].Artists[0].Name)
	suite.Equal(orgSnapshot.Tracks[1].Track.Artists[0].Name, recSnapshot.Tracks[1].Artists[0].Name)
}

func (suite *FavTracksTestSuite) TestFavTracksSnapshotDiff() {
	relSnapshot := suite.snapshots[0]
	req := suite.getRequest("/api/ssfavtracks/diff/{timestamp}")
	req = mux.SetURLVars(req, map[string]string{"timestamp": strconv.FormatInt(relSnapshot.Timestamp.Unix(), 10)})
	req.AddCookie(suite.cookie)

	resp := httptest.NewRecorder()
	suite.handler.ServeHTTP(resp, req)

	suite.NotNil(resp.Body)
	apiResp := suite.checkFavTracksSnapshotDiffAPIResponse(resp.Body.Bytes())
	suite.NotNil(apiResp)

	suite.Nil(apiResp.Results.RemovedTracks)
	suite.NotNil(apiResp.Results.NewTracks)
	suite.Equal(2, len(apiResp.Results.NewTracks))
	suite.Equal(suite.allTracks[2].Track.Name, apiResp.Results.NewTracks[0].Name)
	suite.Equal(suite.allTracks[3].Track.Name, apiResp.Results.NewTracks[1].Name)
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

func (suite *FavTracksTestSuite) checkAllFavTracksSnapshotsAPIResponse(rawResp []byte) *allFavTracksSnapshotsAPIResponse {
	apiResp := &allFavTracksSnapshotsAPIResponse{}
	err := json.Unmarshal(rawResp, apiResp)
	if err != nil {
		suite.FailNowf("fail to unmarshal allFavTracksSnapshotsAPIResponse", "Detals: %s", err.Error())
	}
	suite.Equal(200, apiResp.Status)
	suite.Equal("success", apiResp.Message)
	suite.NotNil(apiResp.Snapshots, "API response tracks must not be nil")
	return apiResp
}

func (suite *FavTracksTestSuite) checkFavTracksSnapshotAPIResponse(rawResp []byte) *favTracksSnapshotAPIresponse {
	apiResp := &favTracksSnapshotAPIresponse{}
	err := json.Unmarshal(rawResp, apiResp)
	if err != nil {
		suite.FailNowf("fail to unmarshal favTracksSnapshotAPIresponse", "Detals: %s", err.Error())
	}
	suite.Equal(200, apiResp.Status)
	suite.Equal("success", apiResp.Message)
	suite.NotNil(apiResp.Snapshot, "API response tracks must not be nil")
	return apiResp
}

func (suite *FavTracksTestSuite) checkFavTracksSnapshotDiffAPIResponse(rawResp []byte) *favTracksSnapshotDiffAPIResponse {
	apiResp := &favTracksSnapshotDiffAPIResponse{}
	err := json.Unmarshal(rawResp, apiResp)
	if err != nil {
		suite.FailNowf("fail to unmarshal favTracksSnapshotAPIresponse", "Detals: %s", err.Error())
	}
	suite.Equal(200, apiResp.Status)
	suite.Equal("success", apiResp.Message)
	suite.NotNil(apiResp.Results, "API response tracks must not be nil")
	return apiResp
}

func (suite *FavTracksTestSuite) fillSnapshotsTestData() {
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
	tr4 := models.SpAddedTrack{
		AddedAt: time.Date(2019, time.August, 15, 10, 0, 0, 0, time.UTC),
		Track: models.SpTrack{
			ID: "track4",
			Artists: []models.SpArtist{
				{
					ID:   "track4id",
					Name: "favTrack4 Artist",
					Type: "track4artistType",
					Href: "dummy href 4",
				},
			},
			Explicit: true,
			Type:     "track 4 type",
			Album: models.SpAlbum{
				ID: "track 4 album ID",
			},
			Name:        "track 4",
			TrackNumber: 9,
		},
	}
	ft1tracks := []models.SpAddedTrack{tr1, tr2}
	ft2tracks := []models.SpAddedTrack{tr3}
	suite.allTracks = []models.SpAddedTrack{tr1, tr2, tr3, tr4}

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

type allFavTracksSnapshotsAPIResponse struct {
	Status    int              `json:"status"`
	Message   string           `json:"message"`
	Snapshots []tracksSnapshot `json:"data"`
}

type favTracksSnapshotAPIresponse struct {
	Status   int            `json:"status"`
	Message  string         `json:"message"`
	Snapshot tracksSnapshot `json:"data"`
}

type favTracksSnapshotDiffAPIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Results struct {
		NewTracks     []track `json:"newTracks"`
		RemovedTracks []track `json:"removedTracks"`
	} `json:"data"`
}

type tracksSnapshot struct {
	Timestamp   int     `json:"timestamp"`
	TracksCount int     `json:"tracks_count"`
	Tracks      []track `json:"tracks"`
}

type track struct {
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
}
