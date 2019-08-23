package db

import "github.com/2beens/spotilizer/models"

type UsersDBTestClient struct {
	testUsers map[string]*models.User
}

func NewUsersDBTest(testUsers []models.User) *UsersDBTestClient {
	dbTestClient := &UsersDBTestClient{}
	dbTestClient.testUsers = make(map[string]*models.User)
	for _, u := range testUsers {
		u := u
		dbTestClient.testUsers[u.Username] = &u
	}
	return dbTestClient
}

func (c UsersDBTestClient) SaveUser(user *models.User) (stored bool) {
	c.testUsers[user.Username] = user
	return true
}

func (c UsersDBTestClient) GetUser(username string) *models.User {
	return c.testUsers[username]
}

func (c UsersDBTestClient) GetAllUsers() []models.User {
	var users []models.User
	for _, user := range c.testUsers {
		users = append(users, *user)
	}
	return users
}
