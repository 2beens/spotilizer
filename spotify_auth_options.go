package main

// SpotifyAuthOptions is: https://developer.spotify.com/documentation/general/guides/authorization-guide/
type SpotifyAuthOptions struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
