package services_test

import (
	"log"
	"testing"

	m "github.com/2beens/spotilizer/models"
	s "github.com/2beens/spotilizer/services"
)

type cookiesDBClientMock struct{}
type usersDBClientMock struct{}

func (uDB usersDBClientMock) SaveUser(user *m.User) (stored bool) {
	return user != nil
}

func (uDB usersDBClientMock) GetUser(username string) *m.User {
	auth := &m.SpotifyAuthOptions{}
	return &m.User{Username: username, Auth: auth}
}

func (uDB usersDBClientMock) GetAllUsers() *[]m.User {
	users := []m.User{}
	users = append(users, *uDB.GetUser("user1"))
	users = append(users, *uDB.GetUser("user2"))
	return &users
}

func (cDB cookiesDBClientMock) SaveCookiesInfo(cookieID2usernameMap map[string]string) {
	log.Println(" > storing cookies data in DB ...")
	// mock
}

func (cDB cookiesDBClientMock) GetCookiesInfo() (cookieID2usernameMap map[string]string) {
	log.Printf(" > cookiesDBClientMock: creating mock cookies ...")
	cookieID2usernameMap = make(map[string]string)
	cookieID2usernameMap["cookieUser1"] = "user1"
	cookieID2usernameMap["cookieUser2"] = "user2"
	return
}

func failNow(t *testing.T, message string) {
	t.Error(message)
	t.FailNow()
}

func TestUserServiceCookies(t *testing.T) {
	log.Println(" > starting test: TestUserService Cookies")
	defer log.Println("--------------------------------------------------------")

	cookiesDBClientMock := &cookiesDBClientMock{}
	usersDBClientMock := &usersDBClientMock{}
	userService := s.NewUserService(cookiesDBClientMock, usersDBClientMock)

	if userService == nil {
		failNow(t, " >>> error, user service is nil")
	}

	user1cookie, err := userService.GetCookieIDByUsername("user1")
	if err != nil {
		failNow(t, " >>> get cookie for user1 error: "+err.Error())
	}
	if user1cookie != "cookieUser1" {
		failNow(t, " >>> user1cookie not equal")
	}
	log.Println(" > user1 cookie OK")

	user1, err := userService.GetUserByCookieID(user1cookie)
	if err != nil {
		failNow(t, " >>> get user by cookie ID error: "+err.Error())
	}
	if user1 == nil {
		failNow(t, " >>> user1 is nil")
	}
	log.Println(" > user1 is OK")

	user2, err := userService.GetUserByCookieID("cookieUser2")
	if err != nil {
		failNow(t, " >>> get user by cookie ID error: "+err.Error())
	}
	if user2 == nil {
		failNow(t, " >>> user2 is nil")
	}
	log.Println(" > user1 is OK")

	user3, err := userService.GetUserByCookieID("cookieUser3")
	if err == nil || user3 != nil {
		failNow(t, " >>> user3 should not exist")
	}
	log.Println(" > user3 not found, OK")

	userService.Add(&m.User{Username: "user3", Auth: &m.SpotifyAuthOptions{}})
	userService.AddUserCookie("cookieUser3", "user3")
	user3, err = userService.GetUserByCookieID("cookieUser3")
	if err != nil {
		failNow(t, " >>> get user by cookie ID error: "+err.Error())
	}
	if user3 == nil {
		failNow(t, " >>> user3 is nil")
	}
	log.Println(" > user3 is OK")

	un1, found := userService.GetUsernameByCookieID("cookieUser1")
	if len(un1) == 0 || !found {
		failNow(t, " >>> cannot find user1 by it's cookie ID")
	}
	log.Println(" > user1 usernname is OK")

	userService.RemoveUserCookie("cookieUser1")
	un1, found = userService.GetUsernameByCookieID("cookieUser1")
	if len(un1) > 0 || found {
		failNow(t, " >>> error, cookie ID for user1 should be removed")
	}
	log.Println(" > user1 cookie removed, OK")
}

func TestUserServiceUsers(t *testing.T) {
	log.Println(" > starting test: TestUserService Users")
	defer log.Println("--------------------------------------------------------")

	cookiesDBClientMock := &cookiesDBClientMock{}
	usersDBClientMock := &usersDBClientMock{}
	userService := s.NewUserService(cookiesDBClientMock, usersDBClientMock)

	user1, _ := userService.Get("user1")
	if user1 == nil {
		failNow(t, " >>> user1 is nil")
	}
	log.Println(" > user1 is OK")

	// non existing user
	testuser, _ := userService.Get("testuser")
	if testuser != nil {
		failNow(t, " >>> error, testuser should be nil")
	}
	log.Println(" > testuser is nil, OK")

	user2, _ := userService.Get("user2")
	if user2 == nil {
		failNow(t, " >>> user2 is nil")
	}
	log.Println(" > user2 is OK")

	// add user, and then try to get it
	user3, err := userService.Get("user3")
	if user3 != nil || err == nil {
		failNow(t, " >>> error, user3 is not nil")
	}
	userService.Add(&m.User{Username: "user3", Auth: &m.SpotifyAuthOptions{}})
	user3, _ = userService.Get("user3")
	if user3 == nil {
		failNow(t, " >>> user3 is nil")
	}
	log.Println(" > user3 is OK")

	found := userService.Exists("user3")
	if !found {
		failNow(t, " >>> error, user3 shoud be found")
	}
	log.Println(" > user3 found, OK")

	found = userService.Exists("user4")
	if found {
		failNow(t, " >>> error, user4 shoud not be found")
	}
	log.Println(" > user4 not found, OK")
}
