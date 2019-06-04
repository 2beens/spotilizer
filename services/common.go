package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	m "github.com/2beens/spotilizer/models"
)

const requestTimeoutSeconds = 30

func getFromSpotify(apiURL string, path string, authOptions *m.SpotifyAuthOptions) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL+path, nil)
	if err != nil {
		log.Printf(" >>> error getting spotify response. details: %v\n", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authOptions.AccessToken)

	// complicating this function on purpose to demonstrante the usage of channels and goroutines
	// through implementing a request timeout mechanism
	timeoutChan := time.After(time.Duration(requestTimeoutSeconds) * time.Second)
	respChannel := make(chan []byte)
	errChannel := make(chan error)

	go func() {
		resp, err := client.Do(req)
		if err != nil {
			errChannel <- err
			return
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			errChannel <- err
			return
		}
		respChannel <- body
	}()

	select {
	case <-timeoutChan:
		return nil, errors.New("timeout occured")
	case body = <-respChannel:
		return body, nil
	case err = <-errChannel:
		return nil, err
	}
}

func getAPIError(body []byte) (spErr m.SpAPIError, isError bool) {
	err := json.Unmarshal(body, &spErr)
	if err == nil && len(spErr.Error.Message) > 0 && spErr.Error.Status > 0 {
		return spErr, true
	}
	return m.SpAPIError{}, false
}
