package services

import (
	"io/ioutil"
	"log"
	"net/http"

	m "github.com/2beens/spotilizer/models"
)

func getFromSpotify(apiURL string, path string, authOptions m.SpotifyAuthOptions) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL+path, nil)
	if err != nil {
		log.Printf(" >>> error getting spotify response. details: %v\n", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authOptions.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf(" >>> error getting current user playlist. details: %v\n", err)
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(" >>> error getting spotify response. details: %v\n", err)
		return nil, err
	}

	// log.Println(string(body))

	return
}
