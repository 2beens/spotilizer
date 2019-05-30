package constants

const (
	IPAddress       = "localhost"
	CookieStateKey  = "spotify_auth_state"
	CookieUserIDKey = "spotilizer-user-id"
	Protocol        = "http"
	Port            = "8080"
	RedisPort       = "6379"
	Permissions     = `
		user-read-private 
		user-read-email 
		user-library-read 
		user-read-birthdate
		playlist-read-private`
)
