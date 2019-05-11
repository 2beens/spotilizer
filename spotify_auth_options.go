package main

import "fmt"

// SpotifyAuthOptions is: https://developer.spotify.com/documentation/general/guides/authorization-guide/
type SpotifyAuthOptions struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (ao SpotifyAuthOptions) String() string {
	return fmt.Sprintf("Spotify Auth Options [tokenType: %s] [scope: %s] [expires in: %v]", ao.TokenType, ao.Scope, ao.ExpiresIn)
}
