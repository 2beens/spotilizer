package services_test

import (
	"log"
	"testing"

	s "github.com/2beens/spotilizer/services"
)

type CookiesDBClient struct{}

func (self CookiesDBClient) SaveCookiesInfo(cookieID2usernameMap map[string]string) {
	log.Println(" > storing cookies data in DB ...")
	// mock
}

func (self CookiesDBClient) GetCookiesInfo() (cookieID2usernameMap map[string]string) {
	cookieID2usernameMap = make(map[string]string)
	cookieID2usernameMap["testuser"] = "testcookie"
	return
}

func TestUserService(t *testing.T) {
	cookiesDBClientMock := &CookiesDBClient{}
	userService := s.NewUserService(cookiesDBClientMock, nil)
	user, _ := userService.Get("testuser")
	if user != nil {
		t.Error(" > user should be nil")
	}
}
