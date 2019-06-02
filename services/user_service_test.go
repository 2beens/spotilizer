package services_test

import (
	"log"
	"testing"

	m "github.com/2beens/spotilizer/models"
	s "github.com/2beens/spotilizer/services"
)

type CookiesDBClientMock struct{}
type UsersDBClientMock struct{}

func (self UsersDBClientMock) SaveUser(user *m.User) (stored bool) {
	if user == nil {
		return false
	}
	return true
}

func (self UsersDBClientMock) GetUser(username string) *m.User {
	auth := &m.SpotifyAuthOptions{}
	return &m.User{Username: username, Auth: auth}
}

func (self UsersDBClientMock) GetAllUsers() *[]m.User {
	users := []m.User{}
	users = append(users, *self.GetUser("user1"))
	users = append(users, *self.GetUser("user2"))
	return &users
}

func (self CookiesDBClientMock) SaveCookiesInfo(cookieID2usernameMap map[string]string) {
	log.Println(" > storing cookies data in DB ...")
	// mock
}

func (self CookiesDBClientMock) GetCookiesInfo() (cookieID2usernameMap map[string]string) {
	cookieID2usernameMap = make(map[string]string)
	cookieID2usernameMap["user1"] = "testcookie"
	return
}

func TestUserService(t *testing.T) {
	t.Log(" > starting test: TestUserService")

	cookiesDBClientMock := &CookiesDBClientMock{}
	usersDBClientMock := &UsersDBClientMock{}
	userService := s.NewUserService(cookiesDBClientMock, usersDBClientMock)

	user1, _ := userService.Get("user1")
	if user1 == nil {
		t.Error(" >>> user1 is nil")
	} else {
		t.Log(" > user1 is OK")
	}

	// non existing user
	testuser, _ := userService.Get("testuser")
	if testuser == nil {
		t.Log(" > testuser is nil, OK")
	} else {
		t.Error(" > testuser should be nil")
	}

	user2, _ := userService.Get("user2")
	if user2 == nil {
		t.Error(" >>> user2 is nil")
	} else {
		t.Log(" > user2 is OK")
	}

	// add user, and then try to get it
	userService.Add(&m.User{Username: "user3", Auth: &m.SpotifyAuthOptions{}})
	user3, _ := userService.Get("user3")
	if user3 == nil {
		t.Error(" >>> user3 is nil")
	} else {
		t.Log(" > user3 is OK")
	}
}
