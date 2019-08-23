package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/2beens/spotilizer/constants"
	"github.com/2beens/spotilizer/models"
	"github.com/2beens/spotilizer/services"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	assert.Equal(t, 123, 123, "they should be equal")
	assert.NotEqual(t, 123, 456, "they should not be equal")
	u := models.User{
		Username: "aa",
		Auth:     nil,
	}
	if assert.NotNil(t, u) {
		assert.Equal(t, "aa", u.Username)
	}
}

func TestServeHttp(t *testing.T) {
	favTracksHandler := getTestFavTracksHandler()

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

	// TODO: assertations
}

func getTestFavTracksHandler() *FavTracksHandler {
	userPlaylistSrv := services.NewUserPlaylistTestService()
	testUserSrv := services.NewUserServiceTest()
	testUser := &models.User{
		Username: "testUser1",
		Auth: &models.SpotifyAuthOptions{
			AccessToken:  "testat",
			RefreshToken: "testrt",
		},
	}
	testUserSrv.Add(testUser)
	testUserSrv.AddUserCookie("cookietu1", testUser.Username)
	return NewFavTracksHandler(testUserSrv, userPlaylistSrv)
}
